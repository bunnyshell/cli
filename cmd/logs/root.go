package logs

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bunnyshell.com/cli/pkg/api/component"
	"bunnyshell.com/cli/pkg/api/environment"
	"bunnyshell.com/cli/pkg/config"
	k8sLogs "bunnyshell.com/cli/pkg/k8s/kubectl/logs"
	k8sWizard "bunnyshell.com/cli/pkg/wizard/k8s"
	"bunnyshell.com/sdk"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	corev1 "k8s.io/api/core/v1"
)

var mainCmd *cobra.Command

type LogsOptions struct {
	// Component selection
	ComponentIDs   []string
	EnvironmentID  string
	ComponentNames []string

	// Pod/Container selection
	Namespace     string
	PodName       string
	Container     string
	AllContainers bool

	// Log filtering (kubectl standard)
	Follow    bool
	Tail      int64
	Since     string
	SinceTime string
	Timestamps bool
	Previous  bool

	// Output options
	Prefix  bool
	NoColor bool

	OverrideClusterServer string
}

func (o *LogsOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	// Component selection
	flags.StringSliceVar(&o.ComponentIDs, "component", o.ComponentIDs, "Component ID(s) (comma-separated)")
	flags.StringVar(&o.EnvironmentID, "environment", o.EnvironmentID, "Environment ID (stream logs from all components)")
	flags.StringSliceVar(&o.ComponentNames, "name", o.ComponentNames, "Filter by component name (requires --environment, repeatable)")

	// Pod/Container selection
	flags.StringVarP(&o.Namespace, "namespace", "n", o.Namespace, "Kubernetes namespace")
	flags.StringVar(&o.PodName, "pod", o.PodName, "Pod name (interactive selection if not specified)")
	flags.StringVarP(&o.Container, "container", "c", o.Container, "Container name (interactive selection if not specified)")
	flags.BoolVar(&o.AllContainers, "all-containers", o.AllContainers, "Stream from all containers in pod")

	// Log filtering
	flags.BoolVarP(&o.Follow, "follow", "f", o.Follow, "Stream logs continuously")
	flags.Int64Var(&o.Tail, "tail", -1, "Show last N lines (default: all)")
	flags.StringVar(&o.Since, "since", o.Since, "Logs newer than duration (e.g. 5s, 2m, 3h)")
	flags.StringVar(&o.SinceTime, "since-time", o.SinceTime, "Logs after timestamp (RFC3339)")
	flags.BoolVar(&o.Timestamps, "timestamps", o.Timestamps, "Include timestamps in output")
	flags.BoolVar(&o.Previous, "previous", o.Previous, "Show logs from previous terminated container")

	// Output options
	flags.BoolVar(&o.Prefix, "prefix", true, "Prefix lines with source (component/pod/container)")
	flags.BoolVar(&o.NoColor, "no-color", o.NoColor, "Disable color-coded prefixes")

	flags.StringVar(&o.OverrideClusterServer, "cluster-server", o.OverrideClusterServer, "Override kubeconfig cluster server")
}

func init() {
	settings := config.GetSettings()

	logsOptions := LogsOptions{}

	mainCmd = &cobra.Command{
		Use:   "logs [flags]",
		Short: "Stream logs from component containers",
		Long: `Stream container application logs (stdout/stderr) from Kubernetes pods.

This command streams logs from one or more components in an environment. When multiple
components are specified, their logs are merged with color-coded prefixes for easy
identification.

Component Selection:
  --component <id>,<id>    Stream logs from specific component(s)
  --environment <env-id>   Stream logs from all components in environment
  --name <name>            Filter by component name (requires --environment, repeatable)

If no component is specified, the component from your current context will be used.`,
		Example: `  # Stream logs from component in context with follow mode
  bns configure set-context --component comp-123
  bns logs --follow --tail 100

  # Stream logs from multiple specific components
  bns logs --component comp-1,comp-2,comp-3 --follow

  # Stream logs from all components in environment
  bns logs --environment env-123 --tail 50

  # Filter by component name
  bns logs --environment env-123 --name api --name worker --follow

  # Specific pod and container
  bns logs --component comp-123 --pod my-pod --container api --follow

  # All containers in component
  bns logs --component comp-123 --all-containers --follow

  # Time-based filtering
  bns logs --component comp-123 --since 5m --timestamps

  # Previous container logs
  bns logs --component comp-123 --previous`,

		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			// Validate flags
			if err := validateFlags(&logsOptions, settings); err != nil {
				return err
			}

			// Resolve component list
			components, err := resolveComponents(&logsOptions, settings)
			if err != nil {
				return err
			}

			if len(components) == 0 {
				return errors.New("no components found")
			}

			// Get environment ID from first component
			envID, err := getEnvironmentID(components[0])
			if err != nil {
				return err
			}

			// Fetch kubeconfig once for all components
			kubeConfigOptions := environment.NewKubeConfigOptions(envID, logsOptions.OverrideClusterServer)
			kubeConfig, err := environment.KubeConfig(kubeConfigOptions)
			if err != nil {
				return err
			}

			// Build stream sources for all components
			sources, err := buildStreamSources(components, kubeConfig, &logsOptions)
			if err != nil {
				return err
			}

			if len(sources) == 0 {
				return errors.New("no log sources found (no running pods)")
			}

			// Get output format from settings
			outputFormat := settings.OutputFormat
			if outputFormat == "" {
				outputFormat = "stylish"
			}

			// Create multiplexer
			mux := k8sLogs.NewMultiplexer(sources, logsOptions.Prefix, logsOptions.NoColor, outputFormat)

			// Setup signal handling for graceful shutdown
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

			go func() {
				<-sigChan
				fmt.Fprintln(os.Stderr, "\nReceived interrupt signal, stopping log streams...")
				mux.Stop()
			}()

			// Start streaming
			if err := mux.Start(); err != nil {
				return err
			}

			// Wait for completion
			errs := mux.Wait()

			// Report any errors
			if len(errs) > 0 {
				fmt.Fprintln(os.Stderr, "\nErrors occurred during log streaming:")
				for _, err := range errs {
					fmt.Fprintf(os.Stderr, "  %v\n", err)
				}
				return errors.New("log streaming completed with errors")
			}

			return nil
		},
	}

	flags := mainCmd.Flags()
	logsOptions.UpdateFlagSet(flags)

	config.MainManager.CommandWithAPI(mainCmd)
}

func validateFlags(opts *LogsOptions, settings *config.Settings) error {
	// Check mutually exclusive flags
	if len(opts.ComponentIDs) > 0 && opts.EnvironmentID != "" {
		return errors.New("--component and --environment are mutually exclusive")
	}

	// --name requires --environment
	if len(opts.ComponentNames) > 0 && opts.EnvironmentID == "" {
		return errors.New("--name requires --environment")
	}

	// Validate that we have at least one component source
	if len(opts.ComponentIDs) == 0 && opts.EnvironmentID == "" && settings.Profile.Context.ServiceComponent == "" {
		return errors.New("component ID required: use --component, --environment, or set in context with 'bns configure set-context --component <id>'")
	}

	return nil
}

func resolveComponents(opts *LogsOptions, settings *config.Settings) ([]interface{}, error) {
	var components []interface{}

	if len(opts.ComponentIDs) > 0 {
		// Specific component IDs
		for _, id := range opts.ComponentIDs {
			itemOpts := component.NewItemOptions(id)
			comp, err := component.Get(itemOpts)
			if err != nil {
				return nil, fmt.Errorf("failed to get component %s: %w", id, err)
			}
			components = append(components, comp)
		}
	} else if opts.EnvironmentID != "" {
		// All components in environment (optionally filtered by name)
		listOpts := component.NewListOptions()
		listOpts.Environment = opts.EnvironmentID

		result, err := component.List(listOpts)
		if err != nil {
			return nil, fmt.Errorf("failed to list components: %w", err)
		}

		var allComponents []sdk.ComponentCollection
		if result.Embedded != nil {
			allComponents = result.Embedded.Item
		}

		// Filter by name if specified
		if len(opts.ComponentNames) > 0 {
			nameMap := make(map[string]bool)
			for _, name := range opts.ComponentNames {
				nameMap[name] = true
			}

			for _, comp := range allComponents {
				if nameMap[comp.GetName()] {
					components = append(components, comp)
				}
			}
		} else {
			for _, comp := range allComponents {
				components = append(components, comp)
			}
		}
	} else {
		// Fall back to context
		componentID := settings.Profile.Context.ServiceComponent
		if componentID == "" {
			return nil, errors.New("no component specified")
		}

		itemOpts := component.NewItemOptions(componentID)
		comp, err := component.Get(itemOpts)
		if err != nil {
			return nil, fmt.Errorf("failed to get component from context: %w", err)
		}
		components = append(components, comp)
	}

	return components, nil
}

// Helper functions to extract component information from different SDK types
func getComponentID(comp interface{}) (string, error) {
	switch c := comp.(type) {
	case *sdk.ComponentItem:
		return c.GetId(), nil
	case sdk.ComponentCollection:
		return c.GetId(), nil
	default:
		return "", fmt.Errorf("unsupported component type: %T", comp)
	}
}

func getComponentName(comp interface{}) (string, error) {
	switch c := comp.(type) {
	case *sdk.ComponentItem:
		return c.GetName(), nil
	case sdk.ComponentCollection:
		return c.GetName(), nil
	default:
		return "", fmt.Errorf("unsupported component type: %T", comp)
	}
}

func getEnvironmentID(comp interface{}) (string, error) {
	switch c := comp.(type) {
	case *sdk.ComponentItem:
		return c.GetEnvironment(), nil
	case sdk.ComponentCollection:
		return c.GetEnvironment(), nil
	default:
		return "", fmt.Errorf("unsupported component type: %T", comp)
	}
}

func buildStreamSources(
	components []interface{},
	kubeConfig *environment.KubeConfigItem,
	opts *LogsOptions,
) ([]*k8sLogs.StreamSource, error) {
	var sources []*k8sLogs.StreamSource

	// Parse log options
	logOpts := parseLogOptions(opts)

	for _, comp := range components {
		// Get component ID and name
		compID, err := getComponentID(comp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to get component ID: %v\n", err)
			continue
		}

		compName, err := getComponentName(comp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to get component name: %v\n", err)
			continue
		}

		// Get pods for this component
		podListOpts := &k8sWizard.PodListOptions{
			Component: compID,
		}

		pods, err := k8sWizard.PodList(podListOpts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to list pods for component %s: %v\n",
				compName, err)
			continue
		}

		if len(pods) == 0 {
			fmt.Fprintf(os.Stderr, "Warning: no pods found for component %s\n", compName)
			continue
		}

		// For each pod, get containers and create stream sources
		for _, pod := range pods {
			podSources, err := buildPodStreamSources(compID, compName, pod, kubeConfig, opts, logOpts)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to build streams for pod %s: %v\n",
					pod.GetName(), err)
				continue
			}

			sources = append(sources, podSources...)
		}
	}

	return sources, nil
}

func buildPodStreamSources(
	compID string,
	compName string,
	pod sdk.ComponentResourceItem,
	kubeConfig *environment.KubeConfigItem,
	opts *LogsOptions,
	logOpts *k8sLogs.Options,
) ([]*k8sLogs.StreamSource, error) {
	var sources []*k8sLogs.StreamSource

	// Filter by pod name if specified
	if opts.PodName != "" && pod.GetName() != opts.PodName {
		return sources, nil
	}

	// Filter by namespace if specified
	if opts.Namespace != "" && pod.GetNamespace() != opts.Namespace {
		return sources, nil
	}

	// Get containers for the pod
	containers, err := getContainersForPod(kubeConfig, pod.GetNamespace(), pod.GetName())
	if err != nil {
		return nil, err
	}

	// Create stream source for each container
	for _, container := range containers {
		// Filter by container name if specified
		if opts.Container != "" && container.Name != opts.Container {
			continue
		}

		// Create log streamer for this container
		containerLogOpts := *logOpts
		containerLogOpts.Namespace = pod.GetNamespace()
		containerLogOpts.PodName = pod.GetName()
		containerLogOpts.Container = container.Name

		streamer, err := k8sLogs.NewLogStreamer(kubeConfig.Bytes, &containerLogOpts)
		if err != nil {
			return nil, fmt.Errorf("failed to create log streamer: %w", err)
		}

		source := &k8sLogs.StreamSource{
			ComponentID:   compID,
			ComponentName: compName,
			Namespace:     pod.GetNamespace(),
			PodName:       pod.GetName(),
			Container:     container.Name,
			Streamer:      streamer,
		}

		sources = append(sources, source)

		// If not streaming all containers, break after first match
		if opts.Container != "" {
			break
		}
	}

	return sources, nil
}

func getContainersForPod(kubeConfig *environment.KubeConfigItem, namespace, podName string) ([]*corev1.Container, error) {
	// Create log streamer temporarily just to get the k8s client
	tempOpts := &k8sLogs.Options{
		Namespace: namespace,
		PodName:   podName,
	}

	streamer, err := k8sLogs.NewLogStreamer(kubeConfig.Bytes, tempOpts)
	if err != nil {
		return nil, err
	}

	// Use the wizard to get containers
	containerListOpts := &k8sWizard.ContainerListOptions{
		Namespace: namespace,
		PodName:   podName,
		Client:    streamer.PodClient,
	}

	containerItems, err := k8sWizard.ContainerList(containerListOpts)
	if err != nil {
		return nil, err
	}

	containers := make([]*corev1.Container, len(containerItems))
	for i, item := range containerItems {
		containers[i] = item.Container
	}

	return containers, nil
}

func parseLogOptions(opts *LogsOptions) *k8sLogs.Options {
	logOpts := &k8sLogs.Options{
		Follow:     opts.Follow,
		Timestamps: opts.Timestamps,
		Previous:   opts.Previous,
	}

	if opts.Tail >= 0 {
		logOpts.Tail = &opts.Tail
	}

	if opts.Since != "" {
		if duration, err := time.ParseDuration(opts.Since); err == nil {
			logOpts.Since = &duration
		}
	}

	if opts.SinceTime != "" {
		if t, err := time.Parse(time.RFC3339, opts.SinceTime); err == nil {
			logOpts.SinceTime = &t
		}
	}

	return logOpts
}

func GetMainCommand() *cobra.Command {
	return mainCmd
}
