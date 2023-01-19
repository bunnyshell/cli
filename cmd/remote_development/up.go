package remote_development

import (
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/environment"
	"bunnyshell.com/cli/pkg/remote_development"
	"github.com/spf13/cobra"
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	var (
		localSyncPath  string
		remoteSyncPath string
		resourcePath   string
		portMappings   []string
	)

	command := &cobra.Command{
		Use: "up",

		ValidArgsFunction: cobra.NoFileCompletions,

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

			environmentResource, err := environment.NewFromWizard(&settings.Profile.Context, resourcePath)
			if err != nil {
				return err
			}

			remoteDevelopment.WithEnvironmentResource(environmentResource)

			// init
			if err = remoteDevelopment.Up(); err != nil {
				return err
			}

			// start
			if err = remoteDevelopment.StartSSHTerminal(); err != nil {
				return err
			}

			return remoteDevelopment.Wait()
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.Organization.GetFlag("organization"))
	flags.AddFlag(options.Project.GetFlag("project"))
	flags.AddFlag(options.Environment.GetFlag("environment"))
	flags.AddFlag(options.ServiceComponent.GetFlag("component"))

	flags.StringVarP(&localSyncPath, "local-sync-path", "l", "", "Local folder path to sync")
	flags.StringVarP(&remoteSyncPath, "remote-sync-path", "r", "", "Remote folder path to sync")
	flags.StringVarP(&resourcePath, "resource", "s", "", "The cluster resource to use (namespace/kind/name format).")
	flags.StringSliceVarP(
		&portMappings,
		"port-forward",
		"f",
		portMappings,
		"Port forward: '8080>3000'\nReverse port forward: '9003<9003'\nComma separated: '8080>3000,9003<9003'",
	)

	mainCmd.AddCommand(command)
}
