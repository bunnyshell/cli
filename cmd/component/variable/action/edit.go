package action

import (
	"io"
	"os"

	"bunnyshell.com/cli/pkg/api/component_variable"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/util"
	"github.com/spf13/cobra"
)

func init() {
	var componentName string
	var componentVariableName string

	settings := config.GetSettings()

	editOptions := component_variable.NewEditOptions("")

	command := &cobra.Command{
		Use: "edit",

		ValidArgsFunction: cobra.NoFileCompletions,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			hasStdin, err := util.IsStdinPresent()
			if err != nil {
				return err
			}

			flags := cmd.Flags()
			if flags.Changed("value") && hasStdin {
				return errMultipleValueInputs
			}

			if flags.Changed("name") && !flags.Changed("component-name") && settings.Profile.Context.ServiceComponent == "" {
				return errComponentRequired
			}

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			flags := cmd.Flags()
			if flags.Changed("value") {
				editOptions.ServiceComponentVariableEditAction.SetValue(flags.Lookup("value").Value.String())
			}

			if componentVariableName != "" {
				componentVariable, err := findComponentVariableByName(componentVariableName, componentName, &settings.Profile)
				if err != nil {
					return err
				}
				editOptions.ID = componentVariable.GetId()
			}

			hasStdin, err := util.IsStdinPresent()
			if err != nil {
				return err
			}

			if hasStdin {
				buf, err := io.ReadAll(os.Stdin)
				if err != nil {
					return err
				}

				editOptions.ServiceComponentVariableEditAction.SetValue(string(buf))
			}

			model, err := component_variable.Edit(editOptions)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			return lib.FormatCommandData(cmd, model)
		},
	}

	flags := command.Flags()

	flags.AddFlag(GetIDOption(&editOptions.ID).GetFlag("id"))
	flags.StringVar(&componentVariableName, "name", componentVariableName, "Component Variable Name")
	command.MarkFlagsMutuallyExclusive("name", "id")

	updateComponentIdentifierFlags(command, &componentName)
	command.MarkFlagsMutuallyExclusive("component-name", "id", "component")

	editOptions.UpdateFlagSet(flags)

	mainCmd.AddCommand(command)
}
