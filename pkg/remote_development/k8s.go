package remote_development

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"bunnyshell.com/cli/pkg/util"

	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	apiMetaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/watch"
	applyCoreV1 "k8s.io/client-go/applyconfigurations/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

const (
	MetadataRemoteDevPrefix    = "remote-dev.bunnyshell.com/"
	MetadataRemoteDevEnabled   = MetadataRemoteDevPrefix + "enabled"
	MetadataRemoteDevComponent = MetadataRemoteDevPrefix + "component"
	MetadataRemoteDevContainer = MetadataRemoteDevPrefix + "container"

	// todo this should match the secret name that remote development start api is looking for
	SecretNamePattern = "%s-remote-dev"
	// todo this should match the pvc name that remote development start api created
	PVCName = "%s-remote-dev"

	SecretCertKeyName           = "cert.pem"
	SecretKeyKeyName            = "key.pem"
	SecretConfigKeyName         = "config.xml"
	SecretAuthorizedKeysKeyName = "authorized_keys"

	PortForwardMethod    = "POST"
	PortForwardInterface = "127.0.0.1"
	SSHRemotePort        = 2222
)

type PortForwardOptions struct {
	Interface string

	RemotePort int
	LocalPort  int

	StopChannel  chan struct{}
	ReadyChannel chan struct{}
}

func NewPortForwardOptions() *PortForwardOptions {
	return &PortForwardOptions{
		Interface: PortForwardInterface,

		RemotePort: SSHRemotePort,
		LocalPort:  0,

		StopChannel:  make(chan struct{}),
		ReadyChannel: make(chan struct{}, 1),
	}
}

type KubernetesClient struct {
	Config     clientcmd.ClientConfig
	RESTConfig *rest.Config
	ClientSet  *kubernetes.Clientset

	SSHPortForwardOptions *PortForwardOptions
}

func NewKubernetesClient(kubeConfigPath string) (*KubernetesClient, error) {
	newKubernetes := new(KubernetesClient)

	kubeconfig, err := os.ReadFile(kubeConfigPath)
	if err != nil {
		return newKubernetes, err
	}

	config, err := clientcmd.NewClientConfigFromBytes(kubeconfig)
	if err != nil {
		return newKubernetes, err
	}

	restConfig, err := config.ClientConfig()
	if err != nil {
		return newKubernetes, err
	}

	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return newKubernetes, err
	}

	newKubernetes.Config = config
	newKubernetes.RESTConfig = restConfig
	newKubernetes.ClientSet = clientset
	newKubernetes.SSHPortForwardOptions = NewPortForwardOptions()

	return newKubernetes, nil
}

func (r *RemoteDevelopment) EnsureRemoteDevK8sSecret() error {
	spinner := util.MakeSpinner(" Prepare Secret")
	spinner.Start()
	defer spinner.Stop()

	remoteSyncthingConfigDir := r.getRemoteSyncthingConfigDir()
	syncthingCertData, err := os.ReadFile(remoteSyncthingConfigDir + "/cert.pem")
	if err != nil {
		return err
	}

	syncthingKeyData, err := os.ReadFile(remoteSyncthingConfigDir + "/key.pem")
	if err != nil {
		return err
	}

	syncthingConfigData, err := os.ReadFile(remoteSyncthingConfigDir + "/config.xml")
	if err != nil {
		return err
	}

	sshPublicKeyData, err := os.ReadFile(r.SSHPublicKeyPath)
	if err != nil {
		return err
	}

	namespace, _, err := r.KubernetesClient.Config.Namespace()
	if err != nil {
		return err
	}

	secretName := fmt.Sprintf(SecretNamePattern, r.ContainerName)

	labels := make(map[string]string)
	labels[MetadataRemoteDevEnabled] = "true"
	labels[MetadataRemoteDevComponent] = r.ComponentName

	secretData := make(map[string][]byte)
	secretData[SecretCertKeyName] = syncthingCertData
	secretData[SecretKeyKeyName] = syncthingKeyData
	secretData[SecretConfigKeyName] = syncthingConfigData
	secretData[SecretAuthorizedKeysKeyName] = sshPublicKeyData

	applySecret := applyCoreV1.Secret(secretName, namespace).WithLabels(labels).WithData(secretData)
	_, err = r.KubernetesClient.ClientSet.CoreV1().Secrets(namespace).Apply(context.TODO(), applySecret, apiMetaV1.ApplyOptions{
		FieldManager: "bunnyshell-cli",
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *RemoteDevelopment) DeleteRemoteDevK8sPVC() error {
	pvcName := fmt.Sprintf(PVCName, r.ComponentName)
	return r.KubernetesClient.DeletePVC(pvcName)
}

func (k *KubernetesClient) DeletePVC(name string) error {
	namespace, _, err := k.Config.Namespace()
	if err != nil {
		return err
	}
	fmt.Println(namespace)

	return k.ClientSet.CoreV1().PersistentVolumeClaims(namespace).Delete(context.TODO(), name, apiMetaV1.DeleteOptions{})
}

func (k *KubernetesClient) GetDeployment(name string) (*appsV1.Deployment, error) {
	deployment := &appsV1.Deployment{}
	namespace, _, err := k.Config.Namespace()
	if err != nil {
		return deployment, err
	}

	return k.ClientSet.AppsV1().Deployments(namespace).Get(context.TODO(), name, apiMetaV1.GetOptions{})
}

func (k *KubernetesClient) GetDeploymentContainers(name string) ([]coreV1.Container, error) {
	deployment, err := k.GetDeployment(name)
	return deployment.Spec.Template.Spec.Containers, err
}

func (k *KubernetesClient) WatchRemoteDevPods(componentName, containerName string) error {
	spinner := util.MakeSpinner(" Waiting for pod to be ready for remote development")
	spinner.Start()
	defer spinner.Stop()

	labelSelector := apiMetaV1.LabelSelector{MatchLabels: map[string]string{
		MetadataRemoteDevEnabled:   "true",
		MetadataRemoteDevComponent: componentName,
	}}

	namespace, _, err := k.Config.Namespace()
	if err != nil {
		return err
	}

	timeout := int64(120)
	listOptions := apiMetaV1.ListOptions{
		LabelSelector:  labels.Set(labelSelector.MatchLabels).String(),
		TimeoutSeconds: &timeout,
	}
	podList, err := k.ClientSet.CoreV1().Pods(namespace).List(context.TODO(), listOptions)
	if err != nil {
		return err
	}
	allRunning := true
	for _, pod := range podList.Items {
		if pod.DeletionTimestamp != nil || pod.Status.Phase != coreV1.PodRunning {
			allRunning = false
			break
		}
	}

	if allRunning {
		return nil
	}

	watcher, err := k.ClientSet.CoreV1().Pods(namespace).Watch(context.TODO(), listOptions)
	if err != nil {
		return err
	}

	defer watcher.Stop()
	for event := range watcher.ResultChan() {
		pod := event.Object.(*coreV1.Pod)
		// ignore terminating pod
		if pod.DeletionTimestamp != nil {
			continue
		}

		if event.Type == watch.Added {
			continue
		}

		for _, containerStatus := range pod.Status.ContainerStatuses {
			if containerStatus.Name == containerName && containerStatus.Ready {
				return nil
			}
		}
	}

	// timeout reached
	return fmt.Errorf("failed to start remote development")
}

func (k *KubernetesClient) GetComponentPod(componentName string) (coreV1.Pod, error) {
	nilPod := coreV1.Pod{}
	labelSelector := apiMetaV1.LabelSelector{MatchLabels: map[string]string{
		MetadataRemoteDevEnabled:   "true",
		MetadataRemoteDevComponent: componentName,
	}}

	namespace, _, err := k.Config.Namespace()
	if err != nil {
		return nilPod, err
	}

	listOptions := apiMetaV1.ListOptions{
		LabelSelector: labels.Set(labelSelector.MatchLabels).String(),
	}
	podList, err := k.ClientSet.CoreV1().Pods(namespace).List(context.TODO(), listOptions)
	if err != nil {
		return nilPod, err
	}

	for _, pod := range podList.Items {
		if pod.DeletionTimestamp == nil && pod.Status.Phase == coreV1.PodRunning {
			return pod, nil
		}
	}

	return nilPod, fmt.Errorf("pod not found for component %v", componentName)
}

func (k *KubernetesClient) GetPortForwardSubresourceURL(pod coreV1.Pod) *url.URL {
	return k.ClientSet.CoreV1().RESTClient().Post().
		Resource("pods").
		Namespace(pod.Namespace).
		Name(pod.Name).
		SubResource("portforward").URL()
}

func (k *KubernetesClient) PortForwardRemoteSSH(componentName string) (*portforward.PortForwarder, error) {
	pod, err := k.GetComponentPod(componentName)
	if err != nil {
		return nil, err
	}

	transport, upgrader, err := spdy.RoundTripperFor(k.RESTConfig)
	if err != nil {
		return nil, err
	}

	url := k.GetPortForwardSubresourceURL(pod)
	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, PortForwardMethod, url)

	k.SSHPortForwardOptions.LocalPort, err = util.GetAvailableEphemeralPort(k.SSHPortForwardOptions.Interface)
	if err != nil {
		return nil, err
	}
	ports := []string{fmt.Sprintf(
		"%d:%d",
		k.SSHPortForwardOptions.LocalPort,
		k.SSHPortForwardOptions.RemotePort,
	)}

	forwarder, err := portforward.NewOnAddresses(
		dialer,
		[]string{k.SSHPortForwardOptions.Interface},
		ports,
		k.SSHPortForwardOptions.StopChannel,
		k.SSHPortForwardOptions.ReadyChannel,
		io.Discard,
		os.Stderr,
	)
	if err != nil {
		return nil, err
	}

	errChan := make(chan error, 1)
	go func() {
		errChan <- forwarder.ForwardPorts()
		close(errChan)
	}()

	select {
	case <-forwarder.Ready:
	case err := <-errChan:
		return nil, err
	}

	return forwarder, nil
}
