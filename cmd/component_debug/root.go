package component_debug

import (
	"bunnyshell.com/cli/pkg/config"
	"github.com/spf13/cobra"
)

var mainCmd = &cobra.Command{
	Use:     "debug",
	Aliases: []string{"debug"},

	Short: "Debug Component",
}

func init() {
	config.MainManager.CommandWithAPI(mainCmd)
}

func GetMainCommand() *cobra.Command {
	return mainCmd
}
