package profile

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"bunnyshell.com/cli/pkg/lib"
)

func init() {
	showConfigCommand := &cobra.Command{
		Use: "list",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := lib.GetConfig()
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			return lib.FormatCommandData(cmd, map[string]interface{}{
				"file": viper.ConfigFileUsed(),
				"data": config.Profiles,
			})
		},
	}

	mainCmd.AddCommand(showConfigCommand)
}
