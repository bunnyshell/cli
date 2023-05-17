package config

import (
	"bunnyshell.com/cli/pkg/util"
	"github.com/spf13/cobra"
)

var mainCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"cfg"},

	Short: "Manage Remote Development config",
	Long:  "Manage Remote Development config",

	PersistentPreRunE: util.PersistentPreRunChain,
}

func GetMainCommand() *cobra.Command {
	return mainCmd
}
