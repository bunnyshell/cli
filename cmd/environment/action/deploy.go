package action

import (
	"fmt"

	"bunnyshell.com/cli/pkg/api/environment"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/progress"
	"github.com/spf13/cobra"
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	monitor := false

	command := &cobra.Command{
		Use: "deploy",

		ValidArgsFunction: cobra.NoFileCompletions,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			if !monitor {
				return nil
			}

			if settings.IsStylish() {
				return nil
			}

			return fmt.Errorf("%w when following deployments", lib.ErrNotStylish)
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			model, err := environment.Deploy(settings.Profile.Context.Environment)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			_ = lib.FormatCommandData(cmd, model)

			if monitor {
				if err = progress.Event(model.GetId(), nil); err != nil {
					return lib.FormatCommandError(cmd, err)
				}
			}

			return nil
		},
	}

	flags := command.Flags()

	idFlag := options.Environment.GetFlag("id")
	flags.AddFlag(idFlag)
	_ = command.MarkFlagRequired(idFlag.Name)

	flags.BoolVar(&monitor, "monitor", monitor, "Monitor deployment pipeline until finish. This will become the default for stylish.")

	mainCmd.AddCommand(command)
}
