package remote_development

import (
	"github.com/spf13/cobra"

	"bunnyshell.com/cli/pkg/environment"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/remote_development"
)

func init() {
	var (
		localSyncPath  string
		remoteSyncPath string
		resourcePath   string
		portMappings   []string
	)

	command := &cobra.Command{
		Use: "up",
		RunE: func(cmd *cobra.Command, args []string) error {
			remoteDevelopment := remote_development.NewRemoteDevelopment()

			if localSyncPath != "" {
				remoteDevelopment.WithLocalSyncPath(localSyncPath)
			}

			if remoteSyncPath != "" {
				remoteDevelopment.WithRemoteSyncPath(remoteSyncPath)
			}

			if len(portMappings) > 0 {
				remoteDevelopment.WithPortMappings(portMappings)
			}

			environmentResource, err := environment.NewFromWizard(&lib.CLIContext.Profile.Context, resourcePath)
			if err != nil {
				return err
			}

			remoteDevelopment.WithEnvironmentResource(environmentResource)

			// init
			if err := remoteDevelopment.Up(); err != nil {
				return err
			}

			// start
			if err := remoteDevelopment.StartSSHTerminal(); err != nil {
				return err
			}

			return remoteDevelopment.Wait()
		},
	}

	command.Flags().StringVar(&lib.CLIContext.Profile.Context.ServiceComponent, "component", "", "Service Component")
	command.Flags().StringVarP(&localSyncPath, "local-sync-path", "l", "", "Local folder path to sync")
	command.Flags().StringVarP(&remoteSyncPath, "remote-sync-path", "r", "", "Remote folder path to sync")
	command.Flags().StringVarP(&resourcePath, "resource", "s", "", "The cluster resource to use (namespace/kind/name format).")
	command.Flags().StringSliceVarP(&portMappings, "portforward", "f", []string{}, "Port forward: '8080>3000'\nReverse port forward: '9003<9003'\nComma separated: '8080>3000,9003<9003'")

	mainCmd.AddCommand(command)
}
