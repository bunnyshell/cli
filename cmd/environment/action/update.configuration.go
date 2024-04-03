package action

import (
	"bunnyshell.com/cli/pkg/api/environment"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	editConfigurationOptions := environment.NewEditConfigurationOptions("")

	command := &cobra.Command{
		Use: "update-configuration",

		ValidArgsFunction: cobra.NoFileCompletions,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return editConfigurationOptions.Validate()
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			editConfigurationOptions.ID = settings.Profile.Context.Environment

			if err := editConfigurationOptions.AttachGenesis(); err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			model, err := environment.EditConfiguration(editConfigurationOptions)
			if err != nil {
				return editConfigurationOptions.HandleError(cmd, err)
			}

			if !editConfigurationOptions.WithDeploy {
				return lib.FormatCommandData(cmd, model)
			}

			deployOptions := &editConfigurationOptions.DeployOptions
			deployOptions.ID = model.GetId()

			return HandleDeploy(cmd, deployOptions, "updated", editConfigurationOptions.K8SIntegration, settings.IsStylish())
		},
	}

	flags := command.Flags()

	idFlag := options.Environment.GetFlag("id")
	flags.AddFlag(idFlag)
	_ = command.MarkFlagRequired(idFlag.Name)

	editConfigurationOptions.UpdateCommandFlags(command)

	mainCmd.AddCommand(command)
}
