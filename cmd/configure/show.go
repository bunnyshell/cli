package configure

import (
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	showConfigCommand := &cobra.Command{
		Use: "show",

		Short: "Show current config",
		Long:  "Show currently used CLI config",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			if config.MainManager.Error != nil {
				return lib.FormatCommandError(cmd, config.MainManager.Error)
			}

			return lib.FormatCommandData(cmd, map[string]interface{}{
				"file": config.GetSettings().ConfigFile,
				"data": config.GetConfig(),
			})
		},
	}

	mainCmd.AddCommand(showConfigCommand)
}
