package action

import (
	"bunnyshell.com/cli/pkg/api/component_variable"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	var componentName string
	var componentVariableName string

	settings := config.GetSettings()

	deleteOptions := component_variable.NewDeleteOptions()

	command := &cobra.Command{
		Use: "delete",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			if componentVariableName != "" {
				componentVariable, err := findComponentVariableByName(componentVariableName, componentName, &settings.Profile)
				if err != nil {
					return err
				}
				deleteOptions.ID = componentVariable.GetId()
			}

			err := component_variable.Delete(deleteOptions)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			cmd.Printf("\nComponent variable %s successfully deleted\n", deleteOptions.ID)

			return nil
		},
	}

	flags := command.Flags()

	flags.AddFlag(GetIDOption(&deleteOptions.ID).GetFlag("id"))
	flags.StringVar(&componentVariableName, "name", componentVariableName, "Component Variable Name")
	command.MarkFlagsMutuallyExclusive("name", "id")

	updateComponentIdentifierFlags(command, &componentName)
	command.MarkFlagsMutuallyExclusive("component-name", "id", "component")

	mainCmd.AddCommand(command)
}
