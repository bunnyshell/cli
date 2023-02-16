package profile

import (
	"bunnyshell.com/cli/pkg/util"
	"github.com/spf13/cobra"
)

var mainCmd = &cobra.Command{
	Use: "profiles",

	Short: "Manage profiles",
	Long:  "Manage CLI Configuration profiles",

	PersistentPreRunE: util.PersistentPreRunChain,
}

func GetMainCommand() *cobra.Command {
	return mainCmd
}
