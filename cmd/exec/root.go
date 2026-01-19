package exec

import (
	"errors"
	"strings"

	"bunnyshell.com/cli/pkg/api/component"
	"bunnyshell.com/cli/pkg/api/environment"
	"bunnyshell.com/cli/pkg/config"
	k8sExec "bunnyshell.com/cli/pkg/k8s/kubectl/exec"
	k8sWizard "bunnyshell.com/cli/pkg/wizard/k8s"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var mainCmd *cobra.Command

type ExecOptions struct {
	ComponentID string

	Namespace string
	PodName   string
	Container string

	TTY   bool
	Stdin bool

	OverrideClusterServer string
}

func (o *ExecOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	flags.StringVarP(&o.Namespace, "namespace", "n", o.Namespace, "Kubernetes namespace")
	flags.StringVar(&o.PodName, "pod", o.PodName, "Pod name (supports namespace/pod-name format)")
	flags.StringVarP(&o.Container, "container", "c", o.Container, "Container name")

	flags.BoolVar(&o.TTY, "tty", o.TTY, "Allocate a pseudo-TTY")
	flags.BoolVar(&o.Stdin, "stdin", o.Stdin, "Pass stdin to the container")

	flags.StringVar(&o.OverrideClusterServer, "cluster-server", o.OverrideClusterServer, "Override kubeconfig cluster server with :port, host:port or scheme://host:port")
}

func (o *ExecOptions) MakeExecOptions(kubeConfig *environment.KubeConfigItem, command []string) *k8sExec.Options {
	return &k8sExec.Options{
		TTY:   o.TTY,
		Stdin: o.Stdin,

		Command: command,

		KubeConfig: kubeConfig.Bytes,
	}
}

func init() {
	settings := config.GetSettings()

	execOptions := ExecOptions{
		OverrideClusterServer: "",
	}

	mainCmd = &cobra.Command{
		Use:   "exec [component-id] [flags] -- COMMAND [args...]",
		Short: "Execute a command in a container",
		Long: `Execute a command in a container of a component's pod.

This command is similar to 'kubectl exec' and 'docker exec'. It allows you to run
arbitrary commands in a component's container.

The component ID can be provided as the first positional argument, or it will use
the component ID from your current context.

The '--' separator is used to separate command flags from the command to execute.
Everything after '--' is passed as the command to run in the container.

If no command is specified, defaults to an interactive shell (/bin/sh) with
--tty and --stdin automatically enabled.`,
		Example: `  # Get an interactive shell (auto-enables --tty --stdin)
  bns exec comp-123

  # Run a single command (no TTY/stdin needed)
  bns exec comp-123 -- ls -la /app

  # Specify pod and container explicitly
  bns exec comp-123 --tty --stdin --pod my-pod-abc -c api -- /bin/bash

  # Pipe local script to remote container
  bns exec comp-123 --stdin -- python3 < local-script.py

  # Use component from context (no ID needed)
  bns configure set-context --component comp-123
  bns exec --tty --stdin

  # Specify namespace
  bns exec comp-123 --tty --stdin -n default --pod my-pod -- sh`,

		Args:              cobra.ArbitraryArgs,
		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			// Parse component ID from positional arg or context
			componentID := ""
			commandArgs := args

			// Check if first arg is component ID (not starting with -)
			if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
				componentID = args[0]
				commandArgs = args[1:]
			} else {
				// Fall back to context
				componentID = settings.Profile.Context.ServiceComponent
			}

			if componentID == "" {
				return errors.New("component ID required: provide as argument or set in context with 'bns configure set-context --component <id>'")
			}

			// Determine the command to execute
			execCommand := commandArgs
			if len(execCommand) == 0 {
				// Default to interactive shell
				execCommand = []string{"/bin/sh"}
				// Auto-enable TTY and stdin for interactive shell if not explicitly set
				if !cmd.Flags().Changed("tty") && !cmd.Flags().Changed("stdin") {
					execOptions.TTY = true
					execOptions.Stdin = true
				}
			}

			// Get component from API
			itemOptions := component.NewItemOptions(componentID)
			componentItem, err := component.Get(itemOptions)
			if err != nil {
				return err
			}

			// Fetch kubeconfig for the environment
			kubeConfigOptions := environment.NewKubeConfigOptions(componentItem.GetEnvironment(), execOptions.OverrideClusterServer)
			kubeConfig, err := environment.KubeConfig(kubeConfigOptions)
			if err != nil {
				return err
			}

			// Create exec command using existing package
			execCommandObj, err := k8sExec.Exec(execOptions.MakeExecOptions(kubeConfig, execCommand))
			if err != nil {
				return err
			}

			// Setup pod list and container list options
			podListOptions := &k8sWizard.PodListOptions{
				Component: componentItem.GetId(),
			}
			containerListOptions := &k8sWizard.ContainerListOptions{
				Client: execCommandObj.PodClient,
			}

			// Ensure pod and container are selected
			if err = ensureContainerSelectedForExec(&execOptions, podListOptions, containerListOptions); err != nil {
				return err
			}

			// Set exec parameters
			execCommandObj.Namespace = execOptions.Namespace
			execCommandObj.PodName = execOptions.PodName
			execCommandObj.ContainerName = execOptions.Container

			// Validate and run
			if err = execCommandObj.Validate(); err != nil {
				return err
			}

			return execCommandObj.Run()
		},
	}

	flags := mainCmd.Flags()

	execOptions.UpdateFlagSet(flags)

	config.MainManager.CommandWithAPI(mainCmd)
}

func ensureContainerSelectedForExec(
	execOptions *ExecOptions,
	podListOptions *k8sWizard.PodListOptions,
	containerListOptions *k8sWizard.ContainerListOptions,
) error {
	// Handle pod selection
	if execOptions.PodName != "" {
		// Support namespace/pod-name format
		parts := strings.Split(execOptions.PodName, "/")
		if len(parts) == 2 {
			execOptions.Namespace = parts[0]
			execOptions.PodName = parts[1]
		} else if len(parts) != 1 {
			return errors.New("invalid pod name format; use 'pod-name' or 'namespace/pod-name'")
		}
	} else {
		// Interactive pod selection
		resource, err := k8sWizard.PodSelect(podListOptions)
		if err != nil {
			return err
		}

		execOptions.PodName = resource.GetName()
		execOptions.Namespace = resource.GetNamespace()
	}

	// Handle container selection
	if execOptions.Container == "" {
		containerListOptions.Namespace = execOptions.Namespace
		containerListOptions.PodName = execOptions.PodName

		container, err := k8sWizard.ContainerSelect(containerListOptions)
		if err != nil {
			return err
		}

		execOptions.Container = container.Name
	}

	return nil
}

func GetMainCommand() *cobra.Command {
	return mainCmd
}
