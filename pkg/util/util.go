package util

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	persistentPreRun = "persistent_pre_run"
	trueStr          = "true"
)

func PersistentPreRunChain(command *cobra.Command, args []string) error {
	return preRunChain(command.Parent(), command, args)
}

func AddGroupedCommands(mainCommand *cobra.Command, group cobra.Group, cmds []*cobra.Command) {
	mainCommand.AddGroup(&group)

	mainCommand.AddCommand(cmds...)

	for _, cmd := range cmds {
		cmd.GroupID = group.ID
	}
}

func AllComandsHelpFlag(command *cobra.Command) {
	flags := command.Flags()

	helpFlagName := "help"
	if flags.Lookup(helpFlagName) == nil {
		flags.BoolP(helpFlagName, "h", false, "Help for "+command.Name())
		_ = flags.SetAnnotation(helpFlagName, cobra.FlagSetByCobraAnnotation, []string{trueStr})
	}

	for _, command := range command.Commands() {
		AllComandsHelpFlag(command)
	}
}

func preRunChain(command *cobra.Command, called *cobra.Command, args []string) error {
	if command == nil {
		return nil
	}

	if command.Annotations == nil {
		command.Annotations = map[string]string{
			persistentPreRun: trueStr,
		}
	} else {
		if command.Annotations[persistentPreRun] == trueStr {
			return nil
		}

		command.Annotations[persistentPreRun] = trueStr
	}

	if err := preRunChain(command.Parent(), called, args); err != nil {
		return err
	}

	if command.PersistentPreRun != nil {
		command.PersistentPreRun(called, args)
	} else if command.PersistentPreRunE != nil {
		if err := command.PersistentPreRunE(called, args); err != nil {
			return fmt.Errorf("%s: %w", command.UseLine(), err)
		}
	}

	return nil
}
