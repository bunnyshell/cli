package port_forward

import (
	"bunnyshell.com/cli/pkg/config"
	"github.com/spf13/cobra"
)

var mainCmd = &cobra.Command{
	Use:     "port-forward",
	Aliases: []string{"pfwd"},

	Short: "Port Forward",
	Long:  "Kubernetes Pod Port Forward",
}

func init() {
	config.MainManager.CommandWithAPI(mainCmd)
}

func GetMainCommand() *cobra.Command {
	return mainCmd
}
