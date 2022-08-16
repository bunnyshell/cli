package configure

import (
	"bunnyshell.com/cli/pkg/lib"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	showConfigCommand := &cobra.Command{
		Use:   "show",
		Short: "Show current config",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := viper.ReadInConfig(); err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			config, err := lib.GetConfig()
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			return lib.FormatCommandData(cmd, map[string]interface{}{
				"file": viper.ConfigFileUsed(),
				"data": config,
			})
		},
	}

	mainCmd.AddCommand(showConfigCommand)
}
