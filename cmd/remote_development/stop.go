package remote_development

import (
	"github.com/spf13/cobra"

	"bunnyshell.com/cli/pkg/lib"
	remoteDevPkg "bunnyshell.com/cli/pkg/remote_development"
)

func init() {
	command := &cobra.Command{
		Use: "stop",
		RunE: func(cmd *cobra.Command, args []string) error {
			organizationId := lib.CLIContext.Profile.Context.Organization
			projectId := lib.CLIContext.Profile.Context.Project
			environmentId := lib.CLIContext.Profile.Context.Environment
			componentId := lib.CLIContext.Profile.Context.ServiceComponent

			remoteDevelopment := remoteDevPkg.NewRemoteDevelopment()

			if err := remoteDevelopment.SelectOrganization(organizationId); err != nil {
				return err
			}

			if err := remoteDevelopment.SelectProject(projectId); err != nil {
				return err
			}

			if err := remoteDevelopment.SelectEnvironment(environmentId); err != nil {
				return err
			}

			if err := remoteDevelopment.SelectComponent(componentId); err != nil {
				return err
			}

			return remoteDevelopment.StopRemoteDevelopment()
		},
	}

	command.Flags().StringVar(&lib.CLIContext.Profile.Context.Organization, "organization", "", "Select Organization")
	command.Flags().StringVar(&lib.CLIContext.Profile.Context.Project, "project", "", "Select Project")
	command.Flags().StringVar(&lib.CLIContext.Profile.Context.Environment, "environment", "", "Select Environment")
	command.Flags().StringVar(&lib.CLIContext.Profile.Context.ServiceComponent, "component", "", "Select Service Component")

	mainCmd.AddCommand(command)
}
