package configure

import (
	"bunnyshell.com/cli/cmd/configure/profile"
	"github.com/spf13/cobra"
)

var mainCmd = &cobra.Command{
	Use:     "configure",
	Aliases: []string{"config"},

	Short: "Configure CLI settings",
}

func init() {
	mainCmd.AddCommand(profile.GetMainCommand())
}

func GetMainCommand() *cobra.Command {
	return mainCmd
}
