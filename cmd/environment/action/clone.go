package action

import (
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
	"github.com/spf13/cobra"
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	cloneName := ""

	command := &cobra.Command{
		Use: "clone",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := lib.GetContext()
			defer cancel()

			request := lib.GetAPI().EnvironmentApi.EnvironmentClone(
				ctx,
				settings.Profile.Context.Environment,
			).EnvironmentCloneAction(
				*sdk.NewEnvironmentCloneAction(cloneName),
			)

			model, resp, err := request.Execute()

			return lib.FormatRequestResult(cmd, model, resp, err)
		},
	}

	flags := command.Flags()
	idFlag := options.Environment.GetFlag("id")

	flags.AddFlag(idFlag)
	_ = command.MarkFlagRequired(idFlag.Name)

	cloneFlagName := "name"
	flags.StringVar(&cloneName, cloneFlagName, cloneName, "Environment Clone Name")
	_ = command.MarkFlagRequired(cloneFlagName)

	mainCmd.AddCommand(command)
}
