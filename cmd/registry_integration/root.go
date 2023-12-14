package registry_integration

import (
	"bunnyshell.com/cli/pkg/config"
	"github.com/spf13/cobra"
)

var mainCmd = &cobra.Command{
	Use:     "container-registries",
	Aliases: []string{"creg"},

	Short: "Container Registry Integrations",
	Long:  "Bunnyshell Container Registry Integrations",
}

var mainGroup = cobra.Group{
	ID:    "container-registries",
	Title: "Commands for Container Registries:",
}

func init() {
	config.MainManager.CommandWithAPI(mainCmd)

	mainCmd.AddGroup(&mainGroup)
}

func GetMainCommand() *cobra.Command {
	return mainCmd
}
