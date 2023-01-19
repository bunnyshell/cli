package variable

import (
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
	"github.com/spf13/cobra"
)

func init() {
	var (
		variableID string
		value      string
	)

	actionGroup := &cobra.Group{
		ID:    "actions",
		Title: "Commands for Environment Variable Actions:",
	}

	command := &cobra.Command{
		Use:     "edit",
		GroupID: actionGroup.ID,

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := lib.GetContext()
			defer cancel()

			request := lib.GetAPI().EnvironmentVariableApi.EnvironmentVariableEdit(
				ctx,
				variableID,
			).EnvironmentVariableEdit(
				*toVariableEdit(value),
			)

			model, resp, err := request.Execute()

			return lib.FormatRequestResult(cmd, model, resp, err)
		},
	}

	flags := command.Flags()

	idFlagName := "id"
	flags.StringVar(&variableID, idFlagName, variableID, "Environment Variable Id")
	_ = command.MarkFlagRequired(idFlagName)

	valueFlagName := "value"
	flags.StringVar(&value, valueFlagName, value, "Environment Variable Value")
	_ = command.MarkFlagRequired(valueFlagName)

	mainCmd.AddGroup(actionGroup)
	mainCmd.AddCommand(command)
}

func toVariableEdit(value string) *sdk.EnvironmentVariableEdit {
	edit := sdk.NewEnvironmentVariableEdit()
	edit.SetValue(value)

	return edit
}
