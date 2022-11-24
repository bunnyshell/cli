package k8s

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"bunnyshell.com/cli/pkg/util"
	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	apiMetaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

const (
	PortForwardMethod = "POST"

	DeploymentKind  = "deployment"
	StatefulSetKind = "statefulset"
	DaemonSetKind   = "daemonset"
)

type KubernetesClient struct {
	kubeConfigPath string

	config     clientcmd.ClientConfig
	restConfig *rest.Config
	clientSet  *kubernetes.Clientset
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

	newKubernetes.kubeConfigPath = kubeConfigPath
	newKubernetes.config = config
	newKubernetes.restConfig = restConfig
	newKubernetes.clientSet = clientset

	return newKubernetes, nil
}

func (k *KubernetesClient) GetPortForwardSubresourceURL(pod *coreV1.Pod) *url.URL {
	return k.clientSet.CoreV1().RESTClient().Post().
		Resource("pods").
		Namespace(pod.Namespace).
		Name(pod.Name).
		SubResource("portforward").URL()
}

func (k *KubernetesClient) ListPods(namespace string, listOptions apiMetaV1.ListOptions) (*coreV1.PodList, error) {
	return k.clientSet.CoreV1().Pods(namespace).List(context.TODO(), listOptions)
}

func (k *KubernetesClient) GetDeployment(namespace, deploymentName string) (*appsV1.Deployment, error) {
	return k.clientSet.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, apiMetaV1.GetOptions{})
}

func (k *KubernetesClient) GetStatefulSet(namespace, name string) (*appsV1.StatefulSet, error) {
	return k.clientSet.AppsV1().StatefulSets(namespace).Get(context.TODO(), name, apiMetaV1.GetOptions{})
}

func (k *KubernetesClient) GetDaemonSet(namespace, name string) (*appsV1.DaemonSet, error) {
	return k.clientSet.AppsV1().DaemonSets(namespace).Get(context.TODO(), name, apiMetaV1.GetOptions{})
}

func (k *KubernetesClient) GetPod(namespace, name string) (*coreV1.Pod, error) {
	return k.clientSet.CoreV1().Pods(namespace).Get(context.TODO(), name, apiMetaV1.GetOptions{})
}

func (k *KubernetesClient) WorkflowPodsList(namespace, kind, name string) (*coreV1.PodList, error) {
	var labelSelector *apiMetaV1.LabelSelector

	switch strings.ToLower(kind) {
	case DeploymentKind:
		deployment, err := k.GetDeployment(namespace, name)
		if err != nil {
			return nil, err
		}

		labelSelector = deployment.Spec.Selector
	case StatefulSetKind:
		statefulset, err := k.GetStatefulSet(namespace, name)
		if err != nil {
			return nil, err
		}

		labelSelector = statefulset.Spec.Selector
	case DaemonSetKind:
		daemonset, err := k.GetDaemonSet(namespace, name)
		if err != nil {
			return nil, err
		}

		labelSelector = daemonset.Spec.Selector
	default:
		return nil, fmt.Errorf("unsupported '%s' resource kind", kind)
	}

	listOptions := apiMetaV1.ListOptions{
		LabelSelector: labels.Set(labelSelector.MatchLabels).String(),
	}
	podsList, err := k.ListPods(namespace, listOptions)
	if err != nil {
		return nil, err
	}

	return podsList, nil
}

func (k *KubernetesClient) PortForward(pod *coreV1.Pod, portForward *PortForward, out, errOut io.Writer) (*portforward.PortForwarder, error) {
	transport, upgrader, err := spdy.RoundTripperFor(k.restConfig)
	if err != nil {
		return nil, err
	}

	url := k.GetPortForwardSubresourceURL(pod)
	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, PortForwardMethod, url)

	if portForward.LocalPort == 0 {
		portForward.LocalPort, err = util.GetAvailableEphemeralPort(portForward.Interface)
		if err != nil {
			return nil, err
		}
	}
	ports := []string{fmt.Sprintf(
		"%d:%d",
		portForward.LocalPort,
		portForward.RemotePort,
	)}

	forwarder, err := portforward.NewOnAddresses(
		dialer,
		[]string{portForward.Interface},
		ports,
		portForward.StopChannel,
		portForward.ReadyChannel,
		out,
		errOut,
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
