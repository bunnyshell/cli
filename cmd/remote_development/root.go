package remote_development

import (
	"bunnyshell.com/cli/pkg/config"
	"github.com/spf13/cobra"
)

var mainCmd = &cobra.Command{
	Use:     "remote-development",
	Aliases: []string{"rdev"},

	Short: "Remote Development",
}

func init() {
	config.MainManager.CommandWithAPI(mainCmd)
}

func GetMainCommand() *cobra.Command {
	return mainCmd
}
