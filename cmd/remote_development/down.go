package remote_development

import (
	"github.com/spf13/cobra"

	"bunnyshell.com/cli/pkg/environment"
	"bunnyshell.com/cli/pkg/lib"
	remoteDevPkg "bunnyshell.com/cli/pkg/remote_development"
)

func init() {
	var resourcePath string

	command := &cobra.Command{
		Use: "down",
		RunE: func(cmd *cobra.Command, args []string) error {
			remoteDevelopment := remoteDevPkg.NewRemoteDevelopment()

			environmentResource, err := environment.NewFromWizard(&lib.CLIContext.Profile.Context, resourcePath)
			if err != nil {
				return err
			}

			remoteDevelopment.WithEnvironmentResource(environmentResource)

			return remoteDevelopment.Down()
		},
	}

	command.Flags().StringVar(&lib.CLIContext.Profile.Context.Organization, "organization", "", "Select Organization")
	command.Flags().StringVar(&lib.CLIContext.Profile.Context.Project, "project", "", "Select Project")
	command.Flags().StringVar(&lib.CLIContext.Profile.Context.Environment, "environment", "", "Select Environment")
	command.Flags().StringVar(&lib.CLIContext.Profile.Context.ServiceComponent, "component", "", "Select Service Component")
	command.Flags().StringVarP(&resourcePath, "resource", "s", "", "The cluster resource to use (namespace/kind/name format).")

	mainCmd.AddCommand(command)
}
