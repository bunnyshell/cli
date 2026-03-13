package organization

import (
	"bunnyshell.com/cli/pkg/config"
	"github.com/spf13/cobra"
)

var mainCmd = &cobra.Command{
	Use:     "organizations",
	Aliases: []string{"organization", "org", "orgs"},

	Short: "Organizations",
	Long:  "Bunnyshell Organizations",
}

func init() {
	config.MainManager.CommandWithAPI(mainCmd)
}

func GetMainCommand() *cobra.Command {
	return mainCmd
}
