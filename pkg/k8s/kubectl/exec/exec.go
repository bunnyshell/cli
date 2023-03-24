package exec

import (
	"os"

	"bunnyshell.com/cli/pkg/build"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	k8sExec "k8s.io/kubectl/pkg/cmd/exec"
	"k8s.io/kubectl/pkg/scheme"
)

type Options struct {
	TTY     bool
	Stdin   bool
	Command []string

	KubeConfig []byte
}

func Exec(options *Options) (*k8sExec.ExecOptions, error) {
	execOptions, err := makeK8sExec(options.KubeConfig)
	if err != nil {
		return nil, err
	}

	execOptions.TTY = options.TTY
	execOptions.Stdin = options.Stdin
	execOptions.Command = options.Command

	return execOptions, nil
}

func makeK8sExec(kubeConfig []byte) (*k8sExec.ExecOptions, error) {
	config, err := makeRestConfig(kubeConfig)
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &k8sExec.ExecOptions{
		StreamOptions: k8sExec.StreamOptions{
			IOStreams: genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr},
		},

		Executor: &k8sExec.DefaultRemoteExecutor{},

		Config:    config,
		PodClient: client.CoreV1(),
	}, nil
}

func makeRestConfig(bytes []byte) (*rest.Config, error) {
	config, err := clientcmd.NewClientConfigFromBytes(bytes)
	if err != nil {
		return nil, err
	}

	restConfig, err := config.ClientConfig()
	if err != nil {
		return nil, err
	}

	setConfigDefaults(restConfig)

	return restConfig, nil
}

func setConfigDefaults(config *rest.Config) *rest.Config {
	config.GroupVersion = &schema.GroupVersion{Group: "", Version: "v1"}
	config.APIPath = "/api"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	if config.UserAgent == "" {
		config.UserAgent = "BunnyCLI+" + build.Version
	}

	return config
}
