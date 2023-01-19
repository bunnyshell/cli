package variable

import (
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	var variableID string

	command := &cobra.Command{
		Use:     "show",
		GroupID: mainGroup.ID,

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := lib.GetContext()
			defer cancel()

			request := lib.GetAPI().EnvironmentVariableApi.EnvironmentVariableView(ctx, variableID)

			model, resp, err := request.Execute()

			return lib.FormatRequestResult(cmd, model, resp, err)
		},
	}

	flags := command.Flags()

	idFlagName := "id"
	flags.StringVar(&variableID, idFlagName, variableID, "Environment Variable Id")
	_ = command.MarkFlagRequired(idFlagName)

	mainCmd.AddCommand(command)
}
