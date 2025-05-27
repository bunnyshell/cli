package action

import (
	"errors"
	"fmt"

	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/environment"
	"bunnyshell.com/cli/pkg/interactive"
	"bunnyshell.com/cli/pkg/port_forward"
	"github.com/spf13/cobra"
)

var errInvalidPortMapping = errors.New("invalid port mapping")

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	var (
		resourcePath          string
		podName               string
		overrideClusterServer string
	)

	command := &cobra.Command{
		Use:     "port-forward mappings...",
		Aliases: []string{"pfwd"},

		Short:   "Starts the port forwarding for the given mappings.",
		Example: "start 8080:80 3306 :9003",

		ValidArgsFunction: cobra.NoFileCompletions,

		Args: func(cmd *cobra.Command, portMappings []string) error {
			if err := cobra.MinimumNArgs(1)(cmd, portMappings); err != nil {
				return err
			}

			for _, portMapping := range portMappings {
				if portMapping == "" || !port_forward.PortMappingExp.MatchString(portMapping) {
					return fmt.Errorf("%w: %s", errInvalidPortMapping, portMapping)
				}
			}

			return nil
		},

		RunE: func(cmd *cobra.Command, portMappings []string) error {
			if podName == "" && config.GetSettings().NonInteractive {
				return interactive.ErrNonInteractive
			}

			portForwardManager := port_forward.NewPortForwardManager()

			portForwardManager.WithPortMappings(portMappings)

			if overrideClusterServer != "" {
				portForwardManager.WithOverrideClusterServer(overrideClusterServer)
			}

			environmentResource, err := environment.NewFromWizard(&settings.Profile.Context, resourcePath)
			if err != nil {
				return err
			}

			_ = portForwardManager.
				WithEnvironmentResource(environmentResource).
				WithWorkspace().
				PrepareKubernetesClient()

			if podName != "" {
				portForwardManager.WithPodName(podName)
			} else {
				if err = portForwardManager.SelectPod(); err != nil {
					return err
				}
			}

			err = portForwardManager.Start()
			if err != nil {
				return err
			}

			_ = portForwardManager.Wait()

			return nil
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.ServiceComponent.GetRequiredFlag("id"))

	flags.StringVarP(&resourcePath, "resource", "s", "", "The cluster resource to use (namespace/kind/name format).")
	flags.StringVar(&podName, "pod", "", "The resource pod to forward ports to.")
	flags.StringVar(&overrideClusterServer, "override-kubeconfig-cluster-server", "", "Override kubeconfig cluster server with :port, host:port or scheme://host:port")

	mainCmd.AddCommand(command)
}
