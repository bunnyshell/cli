package config

import (
	"bunnyshell.com/cli/pkg/config/option"
	"bunnyshell.com/cli/pkg/helper/rdev"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/util"
	"github.com/spf13/cobra"
)

func init() {
	directory := ".bunnyshell"

	command := &cobra.Command{
		Use:     "generate",
		Aliases: []string{"gen"},

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			if err := rdev.Generate(directory); err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			cmd.Println("RemoteDevelopment config generated successfully in " + directory)

			return nil
		},
	}

	flags := command.Flags()

	flags.AddFlag(getDirectoryOption(&directory).GetFlag("directory", util.FlagDirname))

	mainCmd.AddCommand(command)
}

func getDirectoryOption(value *string) *option.String {
	help := "New directory in which to generate the RemoteDev config file"

	option := option.NewStringOption(value)

	option.AddFlagWithExtraHelp("directory", "Directory to generate the config in", help)

	return option
}
