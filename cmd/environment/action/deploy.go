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

	deployOptions := environment.NewDeployOptions("")

	command := &cobra.Command{
		Use: "deploy",

		ValidArgsFunction: cobra.NoFileCompletions,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return validateActionOptions(&deployOptions.ActionOptions)
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			deployOptions.ID = settings.Profile.Context.Environment

			event, err := environment.Deploy(deployOptions)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			if deployOptions.WithoutPipeline {
				return lib.FormatCommandData(cmd, event)
			}

			if err = processEventPipeline(cmd, event, "deploy"); err != nil {
				cmd.Printf("\nEnvironment %s deployment failed\n", deployOptions.ID)

				return err
			}

			cmd.Printf("\nEnvironment %s successfully deployed\n", deployOptions.ID)

			return showEnvironmentEndpoints(cmd, deployOptions.ID)
		},
	}

	flags := command.Flags()

	idFlag := options.Environment.GetFlag("id")
	flags.AddFlag(idFlag)
	_ = command.MarkFlagRequired(idFlag.Name)

	deployOptions.UpdateFlagSet(flags)

	mainCmd.AddCommand(command)
}
