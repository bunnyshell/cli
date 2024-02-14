package secret

import (
	"bunnyshell.com/cli/pkg/config"
	"github.com/spf13/cobra"
)

var mainCmd = &cobra.Command{
	Use:     "secrets",
	Aliases: []string{"sec"},

	Short: "Secrets",
	Long:  "Bunnyshell Secrets",
}

func init() {
	config.MainManager.CommandWithAPI(mainCmd)
}

func GetMainCommand() *cobra.Command {
	return mainCmd
}
