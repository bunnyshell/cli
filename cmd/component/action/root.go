package action

import (
	"github.com/spf13/cobra"
)

var mainCmd = &cobra.Command{}

func GetMainCommand() *cobra.Command {
	return mainCmd
}
