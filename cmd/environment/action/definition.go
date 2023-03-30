package action

import (
	"bunnyshell.com/cli/pkg/api/environment"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	definitionOptions := environment.NewDefinitionOptions("")

	command := &cobra.Command{
		Use:     "definition",
		Aliases: []string{"def"},

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			definitionOptions.ID = settings.Profile.Context.Environment

			definition, err := environment.Definition(definitionOptions)
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

	flags.AddFlag(options.Environment.GetRequiredFlag("id"))

	mainCmd.AddCommand(command)
}
