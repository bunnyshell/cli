package template

import (
	"bunnyshell.com/cli/pkg/api/template"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	settings := config.GetSettings()

	definitionOptions := template.NewDefinitionOptions("")

	command := &cobra.Command{
		Use:     "definition",
		Aliases: []string{"def"},
		GroupID: mainGroup.ID,

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			definition, err := template.Definition(definitionOptions)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			if settings.OutputFormat == "json" {
				return lib.FormatCommandData(cmd, definition.Data)
			}

			cmd.Println(string(definition.Bytes))

			return nil
		},
	}

	flags := command.Flags()

	flags.AddFlag(getIDOption(&definitionOptions.ID).GetRequiredFlag("id"))

	definitionOptions.UpdateFlagSet(flags)

	mainCmd.AddCommand(command)
}
