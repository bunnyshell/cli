package util

import "github.com/spf13/cobra"

func AddGroupedCommands(mainCommand *cobra.Command, group cobra.Group, cmds []*cobra.Command) {
	mainCommand.AddGroup(&group)

	mainCommand.AddCommand(cmds...)

	for _, cmd := range cmds {
		cmd.GroupID = group.ID
	}
}
