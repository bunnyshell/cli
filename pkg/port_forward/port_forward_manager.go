package port_forward

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"bunnyshell.com/cli/pkg/environment"
	"bunnyshell.com/cli/pkg/interactive"
	"bunnyshell.com/cli/pkg/k8s"
	"bunnyshell.com/cli/pkg/port_forward/workspace"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/portforward"
)

const (
	PortForwardDefaultInterface = "127.0.0.1"
)

var (
	ErrNoPods = fmt.Errorf("the selected resource has no pods")

	TerminationSignals = []os.Signal{
		syscall.SIGINT,
		syscall.SIGTERM,
	}
)

type PortForwardManager struct {
	environmentResource *environment.EnvironmentResource
	pod                 *v1.Pod

	workspace *workspace.Workspace

	kubernetesClient      *k8s.KubernetesClient
	overrideClusterServer string

	portForwards   []*k8s.PortForward
	portForwarders []*portforward.PortForwarder
}

func NewPortForwardManager() *PortForwardManager {
	portForwardManager := &PortForwardManager{
		environmentResource:   environment.NewEnvironmentResource(),
		overrideClusterServer: "",
	}

	return portForwardManager
}

func (m *PortForwardManager) WithEnvironmentResource(environmentResource *environment.EnvironmentResource) *PortForwardManager {
	m.environmentResource = environmentResource

	return m
}

func (m *PortForwardManager) WithWorkspace() *PortForwardManager {
	m.workspace = workspace.NewWorkspace(m.environmentResource.Environment.GetId())

	return m
}

func (m *PortForwardManager) WithOverrideClusterServer(overrideClusterServer string) *PortForwardManager {
	m.overrideClusterServer = overrideClusterServer

	return m
}

func (m *PortForwardManager) WithKubernetesClient(kubernetesClient *k8s.KubernetesClient) *PortForwardManager {
	m.kubernetesClient = kubernetesClient

	return m
}

func (m *PortForwardManager) WithPod(pod *v1.Pod) *PortForwardManager {
	m.pod = pod

	return m
}

func (m *PortForwardManager) WithPodName(podName string) *PortForwardManager {
	pod, err := m.kubernetesClient.GetPod(m.environmentResource.ComponentResource.GetNamespace(), podName)
	if err != nil {
		panic(fmt.Errorf("invalid pod name: %s", podName))
	}

	return m.WithPod(pod)
}

func (m *PortForwardManager) WithPortMappings(portMappings []string) *PortForwardManager {
	for _, portMapping := range portMappings {
		match := PortMappingExp.FindStringSubmatch(portMapping)
		if match == nil {
			panic(fmt.Errorf("invalid port mapping: %s", portMapping))
		}

		var localPort int
		var err error

		if match[1] != "" {
			localPort, err = strconv.Atoi(match[1])
			if err != nil {
				panic(fmt.Errorf("invalid port mapping: %s", portMapping))
			}
		} else {
			// We will assign a random ephemeral port
			localPort = 0
		}

		remotePort := localPort
		if match[3] != "" {
			remotePort, err = strconv.Atoi(match[3])
			if err != nil {
				panic(fmt.Errorf("invalid port mapping: %s", portMapping))
			}
		}

		// We should not get this because of the regex, but let's make sure
		if remotePort == 0 {
			panic(fmt.Errorf("invalid port mapping: %s", portMapping))
		}

		m.portForwards = append(m.portForwards, k8s.NewPortForward(PortForwardDefaultInterface, localPort, remotePort))
	}

	return m
}

func (m *PortForwardManager) SelectPod() error {
	componentResource := m.environmentResource.ComponentResource

	// Fetch the resource list of pods
	podsList, err := m.kubernetesClient.WorkflowPodsList(componentResource.GetNamespace(), componentResource.GetKind(), componentResource.GetName())
	if err != nil {
		return err
	}

	if len(podsList.Items) == 0 {
		return ErrNoPods
	}

	if len(podsList.Items) == 1 {
		m.WithPod(&podsList.Items[0])

		return nil
	}

	podNames := []string{}
	podNamesMap := map[string]*v1.Pod{}

	for _, podItem := range podsList.Items {
		pod := podItem
		podNames = append(podNames, pod.Name)
		podNamesMap[pod.Name] = &pod
	}

	_, podName, err := interactive.Choose("Select pod", podNames)
	if err != nil {
		return err
	}

	m.WithPod(podNamesMap[podName])

	return nil
}

func (m *PortForwardManager) PrepareKubernetesClient() error {
	kubeConfigFile, err := m.workspace.DownloadKubeConfig(m.overrideClusterServer)
	if err != nil {
		return err
	}

	kubernetesClient, err := k8s.NewKubernetesClient(kubeConfigFile)
	if err != nil {
		return err
	}

	m.WithKubernetesClient(kubernetesClient)

	return nil
}

func (m *PortForwardManager) Start() error {
	fmt.Printf("Forwarding ports to pod %s/%s...\n\n", m.environmentResource.ComponentResource.GetNamespace(), m.pod.Name)

	for _, portForward := range m.portForwards {
		forwarder, err := m.kubernetesClient.PortForward(m.pod, portForward, os.Stdout, os.Stderr)
		if err != nil {
			return err
		}

		forwardedPorts, err := forwarder.GetPorts()
		if err != nil {
			return err
		}

		if len(forwardedPorts) == 0 {
			return fmt.Errorf("could not create port forward for local port %d to remote port %d", portForward.LocalPort, portForward.RemotePort)
		}

		m.portForwarders = append(m.portForwarders, forwarder)
	}

	return nil
}

func (m *PortForwardManager) Wait() error {
	// exit on cli signal interrupt
	signalTermination := make(chan os.Signal, 1)
	signal.Notify(signalTermination, TerminationSignals...)
	defer signal.Stop(signalTermination)

	sig := <-signalTermination

	m.Close()

	return fmt.Errorf("terminated by signal: %s", sig)
}

func (m *PortForwardManager) Close() {
	// Close the port forwarders
	for _, portForwarder := range m.portForwarders {
		portForwarder.Close()
	}

	// Close the PortForward channels
	for _, portForward := range m.portForwards {
		portForward.Close()
	}
}
