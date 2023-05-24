package template

import (
	"bunnyshell.com/cli/pkg/config/option"
	"bunnyshell.com/cli/pkg/helper/template"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/util"
	"github.com/spf13/cobra"
)

func init() {
	directory := ""

	command := &cobra.Command{
		Use:     "generate",
		Aliases: []string{"gen"},
		GroupID: mainGroup.ID,

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			if err := template.Generate(directory); err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			cmd.Println("Template generated successfully in " + directory)

			return nil
		},
	}

	flags := command.Flags()

	flags.AddFlag(getDirectoryOption(&directory).GetFlag("directory", util.FlagRequired, util.FlagDirname))

	mainCmd.AddCommand(command)
}

func getDirectoryOption(value *string) *option.String {
	help := "New directory in which to generate the template"

	option := option.NewStringOption(value)

	option.AddFlagWithExtraHelp("directory", "Directory to generate the template in", help)

	return option
}
