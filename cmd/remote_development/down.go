package remote_development

import (
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/environment"
	remoteDevPkg "bunnyshell.com/cli/pkg/remote_development"
	"github.com/spf13/cobra"
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	var resourcePath string

	command := &cobra.Command{
		Use: "down",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			remoteDevelopment := remoteDevPkg.NewRemoteDevelopment()

			environmentResource, err := environment.NewFromWizard(&settings.Profile.Context, resourcePath)
			if err != nil {
				return err
			}

			remoteDevelopment.WithEnvironmentResource(environmentResource)

			return remoteDevelopment.Down()
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.Organization.GetFlag("organization"))
	flags.AddFlag(options.Project.GetFlag("project"))
	flags.AddFlag(options.Environment.GetFlag("environment"))
	flags.AddFlag(options.ServiceComponent.GetFlag("component"))

	flags.StringVarP(&resourcePath, "resource", "s", "", "The cluster resource to use (namespace/kind/name format).")

	mainCmd.AddCommand(command)
}
