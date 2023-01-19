package environment

import (
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	command := &cobra.Command{
		Use:     "show",
		GroupID: mainGroup.ID,

		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := lib.GetContext()
			defer cancel()

			request := lib.GetAPI().EnvironmentApi.EnvironmentView(ctx, settings.Profile.Context.Environment)

			model, resp, err := request.Execute()

			return lib.FormatRequestResult(cmd, model, resp, err)
		},
	}

	flags := command.Flags()

	idFlag := options.Environment.GetFlag("id")
	flags.AddFlag(idFlag)
	_ = command.MarkFlagRequired(idFlag.Name)

	mainCmd.AddCommand(command)
}
