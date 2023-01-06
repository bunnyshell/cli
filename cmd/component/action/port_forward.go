package action

import (
	"fmt"

	"bunnyshell.com/cli/pkg/environment"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/port_forward"
	"github.com/spf13/cobra"
)

func init() {
	var (
		resourcePath string
		podName      string
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
					return fmt.Errorf("invalid port mapping: %s", portMapping)
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, portMappings []string) error {
			portForwardManager := port_forward.NewPortForwardManager()

			portForwardManager.WithPortMappings(portMappings)

			environmentResource, err := environment.NewFromWizard(&lib.CLIContext.Profile.Context, resourcePath)
			if err != nil {
				return err
			}

			portForwardManager.
				WithEnvironmentResource(environmentResource).
				PrepareKubernetesClient()

			if podName != "" {
				portForwardManager.WithPodName(podName)
			} else {
				_ = portForwardManager.SelectPod()
			}

			err = portForwardManager.Start()
			if err != nil {
				return err
			}

			_ = portForwardManager.Wait()

			return nil
		},
	}

	command.Flags().StringVar(&lib.CLIContext.Profile.Context.ServiceComponent, "component", "", "Service Component")
	command.Flags().StringVarP(&resourcePath, "resource", "s", "", "The cluster resource to use (namespace/kind/name format).")
	command.Flags().StringVar(&podName, "pod", "", "The resource pod to forward ports to.")

	mainCmd.AddCommand(command)
}
