package repository

import (
	"fmt"

	"bunnyshell.com/cli/pkg/api/template/repository"
	"bunnyshell.com/cli/pkg/build"
	"bunnyshell.com/cli/pkg/config/option"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	itemOptions := repository.NewItemOptions("")

	command := &cobra.Command{
		Use: "show",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			model, err := repository.Get(itemOptions)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			return lib.FormatCommandData(cmd, model)
		},
	}

	flags := command.Flags()

	flags.AddFlag(getIDOption(&itemOptions.ID).GetRequiredFlag("id"))

	mainCmd.AddCommand(command)
}

func getIDOption(value *string) *option.String {
	help := fmt.Sprintf(
		`Find available TemplateRepositories with "%s templates repository list"`,
		build.Name,
	)

	option := option.NewStringOption(value)

	option.AddFlagWithExtraHelp("id", "TemplateRepository ID", help)

	return option
}
