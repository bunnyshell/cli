package template

import (
	"bunnyshell.com/cli/pkg/helper/template"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	directory := "."

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

	flags.StringVar(&directory, "directory", directory, "Directory to generate the template in")
	_ = command.MarkFlagRequired("directory")
	_ = command.MarkFlagDirname("directory")

	mainCmd.AddCommand(command)
}
