package pipeline

import (
	"fmt"

	"bunnyshell.com/cli/pkg/build"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/config/option"
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

func getIDOption(value *string) *option.String {
	help := fmt.Sprintf(
		`Find available Pipelines with "%s pipeline list"`,
		build.Name,
	)

	option := option.NewStringOption(value)

	option.AddFlagWithExtraHelp("id", "Pipeline ID", help)

	return option
}
