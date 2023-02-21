package k8sIntegration

import (
	"bunnyshell.com/cli/pkg/config"
	"github.com/spf13/cobra"
)

var mainCmd = &cobra.Command{
	Use:     "k8s-clusters",
	Aliases: []string{"k8s"},

	Short: "Kubernetes Cluster Integrations",
	Long:  "Bunnyshell Kubernetes Cluster Integrations",
}

var mainGroup = cobra.Group{
	ID:    "k8s-clusters",
	Title: "Commands for Kubernetes Integrations:",
}

func init() {
	config.MainManager.CommandWithAPI(mainCmd)

	mainCmd.AddGroup(&mainGroup)
}

func GetMainCommand() *cobra.Command {
	return mainCmd
}
