package variable

import (
	"bunnyshell.com/cli/pkg/config/option"
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

			request := lib.GetAPI().EnvironmentVariableAPI.EnvironmentVariableEdit(
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

	flags.AddFlag(getIDOption(&variableID).GetRequiredFlag("id"))
	flags.AddFlag(getEditValueOption(&value).GetRequiredFlag("value"))

	mainCmd.AddGroup(actionGroup)
	mainCmd.AddCommand(command)
}

func getEditValueOption(value *string) *option.String {
	help := "Update the value of an environment variable. A deployment will be required for the updates to take effect."

	idOption := option.NewStringOption(value)

	idOption.AddFlagWithExtraHelp("value", "Environment Variable Value", help)

	return idOption
}

func toVariableEdit(value string) *sdk.EnvironmentVariableEdit {
	edit := sdk.NewEnvironmentVariableEdit()
	edit.SetValue(value)

	return edit
}
