package remote_development

import (
	"os"
	"os/signal"

	"github.com/spf13/cobra"

	"bunnyshell.com/cli/pkg/lib"
	remoteDevPkg "bunnyshell.com/cli/pkg/remote_development"
)

func init() {
	var localSyncPath string

	command := &cobra.Command{
		Use:          "start",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			organizationId := lib.CLIContext.Profile.Context.Organization
			projectId := lib.CLIContext.Profile.Context.Project
			environmentId := lib.CLIContext.Profile.Context.Environment
			componentId := lib.CLIContext.Profile.Context.ServiceComponent

			remoteDevelopment := remoteDevPkg.NewRemoteDevelopment()

			if err := remoteDevelopment.EnsureSSHKeys(); err != nil {
				return err
			}

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

			if err := remoteDevelopment.EnsureComponentFolder(); err != nil {
				return err
			}

			if err := remoteDevelopment.EnsureEnvironmentKubeConfig(); err != nil {
				return err
			}

			if err := remoteDevelopment.SelectContainer(); err != nil {
				return err
			}

			if err := remoteDevelopment.SelectLocalSyncFolder(localSyncPath); err != nil {
				return err
			}

			if err := remoteDevelopment.PrepareSyncthing(); err != nil {
				return err
			}

			if err := remoteDevelopment.EnsureRemoteDevK8sSecret(); err != nil {
				return err
			}

			if err := remoteDevelopment.StartRemoteDevelopment(); err != nil {
				return err
			}

			// close channels on cli signal interrupt
			signals := make(chan os.Signal, 1)
			signal.Notify(signals, os.Interrupt)
			defer signal.Stop(signals)
			go remoteDevelopment.CloseOnSignal(signals)

			remoteDevelopment.Wait()
			return nil
		},
	}

	command.Flags().StringVar(&lib.CLIContext.Profile.Context.Organization, "organization", "", "Organization")
	command.Flags().StringVar(&lib.CLIContext.Profile.Context.Project, "project", "", "Project")
	command.Flags().StringVar(&lib.CLIContext.Profile.Context.Environment, "environment", "", "Environment")
	command.Flags().StringVar(&lib.CLIContext.Profile.Context.ServiceComponent, "component", "", "Service Component")
	command.Flags().StringVar(&localSyncPath, "sync-path", "", "Local folder to sync with remote")

	mainCmd.AddCommand(command)
}
