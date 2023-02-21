package pipeline

import (
	"bunnyshell.com/cli/pkg/config"
	"github.com/spf13/cobra"
)

var mainCmd = &cobra.Command{
	Use:     "pipeline",
	Aliases: []string{"pipe"},

	Short: "Pipeline",
	Long:  "Bunnyshell Pipeline",
}

func init() {
	config.MainManager.CommandWithAPI(mainCmd)
}

func GetMainCommand() *cobra.Command {
	return mainCmd
}
