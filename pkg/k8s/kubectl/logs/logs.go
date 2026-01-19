package logs

import (
	"context"
	"io"
	"time"

	"bunnyshell.com/cli/pkg/build"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubectl/pkg/scheme"
)

type Options struct {
	// Pod/Container selection
	Namespace string
	PodName   string
	Container string

	// Log filtering (kubectl standard)
	Follow     bool
	Tail       *int64
	Since      *time.Duration
	SinceTime  *time.Time
	Timestamps bool
	Previous   bool
}

type LogStreamer struct {
	Config    *rest.Config
	PodClient v1.PodsGetter
	Namespace string
	PodName   string
	Container string
	Options   *Options
	Stream    io.ReadCloser
}

// NewLogStreamer creates a new log streamer for a specific pod/container
func NewLogStreamer(kubeConfig []byte, options *Options) (*LogStreamer, error) {
	config, err := makeRestConfig(kubeConfig)
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &LogStreamer{
		Config:    config,
		PodClient: client.CoreV1(),
		Namespace: options.Namespace,
		PodName:   options.PodName,
		Container: options.Container,
		Options:   options,
	}, nil
}

// Start begins streaming logs from the pod/container
func (ls *LogStreamer) Start(ctx context.Context) error {
	podLogOptions := &corev1.PodLogOptions{
		Container:  ls.Container,
		Follow:     ls.Options.Follow,
		Timestamps: ls.Options.Timestamps,
		Previous:   ls.Options.Previous,
	}

	if ls.Options.Tail != nil {
		podLogOptions.TailLines = ls.Options.Tail
	}

	if ls.Options.Since != nil {
		seconds := int64(ls.Options.Since.Seconds())
		podLogOptions.SinceSeconds = &seconds
	}

	if ls.Options.SinceTime != nil {
		metaTime := metav1.NewTime(*ls.Options.SinceTime)
		podLogOptions.SinceTime = &metaTime
	}

	req := ls.PodClient.Pods(ls.Namespace).GetLogs(ls.PodName, podLogOptions)

	stream, err := req.Stream(ctx)
	if err != nil {
		return err
	}

	ls.Stream = stream
	return nil
}

// Read reads from the log stream
func (ls *LogStreamer) Read(p []byte) (int, error) {
	if ls.Stream == nil {
		return 0, io.EOF
	}
	return ls.Stream.Read(p)
}

// Close closes the log stream
func (ls *LogStreamer) Close() error {
	if ls.Stream != nil {
		return ls.Stream.Close()
	}
	return nil
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
