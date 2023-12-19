package action

import (
	"fmt"

	"bunnyshell.com/cli/cmd/environment/action"
	"bunnyshell.com/cli/pkg/api/common"
	"bunnyshell.com/cli/pkg/api/component"
	"bunnyshell.com/cli/pkg/api/environment"
	"bunnyshell.com/cli/pkg/build"
	"bunnyshell.com/cli/pkg/config"
	githelper "bunnyshell.com/cli/pkg/helper/git"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

type EditComponentData struct {
	common.ItemOptions

	K8SIntegration string

	GitTarget string

	WithDeploy bool
}

var commandExample = fmt.Sprintf(`This command updates the Git details for a component

You can update both the repository and the Git branch / ref:
%[1]s%[2]s components update --id dMVwZO5jGN --git-target https://github.com/my-fork/my-repo@my-main

You can update only the Git branch / ref:
%[1]s%[2]s components update --id dMVwZO5jGN --git-target @fix/bug-1234
.`, "\t", build.Name)

func init() {
	settings := config.GetSettings()
	options := config.GetOptions()

	editComponentsData := &EditComponentData{
		ItemOptions: *common.NewItemOptions(""),
	}

	command := &cobra.Command{
		Use: "update",

		Example: commandExample,

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			editComponentsData.ID = settings.Profile.Context.ServiceComponent

			gitRepository, gitRef, err := githelper.ParseGitSec(editComponentsData.GitTarget)
			if err != nil {
				return lib.FormatCommandError(cmd, fmt.Errorf("invalid git spec for %s: %w", editComponentsData.GitTarget, err))
			}

			model, err := component.Get(&editComponentsData.ItemOptions)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			if settings.IsStylish() {
				cmd.Printf(`Updating component "%s" (%s)%s`, editComponentsData.ID, model.GetName(), "\n\n")
			}

			editOptions := environment.NewEditComponentOptions()
			editOptions.ID = model.GetEnvironment()
			editOptions.Component = model.GetName()
			editOptions.TargetRepository = gitRepository
			editOptions.TargetBranch = gitRef

			_, err = environment.EditComponents(editOptions)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			if settings.IsStylish() {
				cmd.Printf("Successfully updated component...%s", "\n\n")
			}

			if !editComponentsData.WithDeploy {
				return lib.FormatCommandData(cmd, model)
			}

			deployOptions := &editOptions.DeployOptions
			deployOptions.ID = model.GetEnvironment()

			if err = action.HandleDeploy(cmd, deployOptions, "updated", editComponentsData.K8SIntegration); err != nil {
				return err
			}

			return lib.FormatCommandData(cmd, model)
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.ServiceComponent.GetRequiredFlag("id"))

	flags.StringVar(&editComponentsData.GitTarget, "git-target", editComponentsData.GitTarget, "Target git spec (e.g. https://github.com/fork/templates@main)")

	flags.BoolVar(&editComponentsData.WithDeploy, "deploy", editComponentsData.WithDeploy, "Deploy the environment after update")
	flags.StringVar(&editComponentsData.K8SIntegration, "k8s", editComponentsData.K8SIntegration, "Set Kubernetes integration for the environment (if not set)")

	mainCmd.AddCommand(command)
}
