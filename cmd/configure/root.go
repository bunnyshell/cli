package configure

import (
	"github.com/spf13/cobra"

	"bunnyshell.com/cli/cmd/configure/profile"
	"bunnyshell.com/cli/pkg/lib"
)

var mainCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure CLI settings",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		lib.MakeDefaultContext()
		return nil
	},
}

func init() {
	mainCmd.AddCommand(profile.GetMainCommand())
}

func GetMainCommand() *cobra.Command {
	return mainCmd
}
