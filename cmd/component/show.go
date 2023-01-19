package component

import (
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	command := &cobra.Command{
		Use: "show",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := lib.GetContext()
			defer cancel()

			request := lib.GetAPI().ComponentApi.ComponentView(ctx, settings.Profile.Context.ServiceComponent)

			model, resp, err := request.Execute()

			return lib.FormatRequestResult(cmd, model, resp, err)
		},
	}

	flags := command.Flags()

	idFlag := options.ServiceComponent.GetFlag("id")
	flags.AddFlag(idFlag)
	_ = command.MarkFlagRequired(idFlag.Name)

	mainCmd.AddCommand(command)
}
