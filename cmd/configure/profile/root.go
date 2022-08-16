package profile

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var mainCmd = &cobra.Command{
	Use:   "profiles",
	Short: "Manage profiles",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if parent := cmd.Parent(); parent != nil {
			if parent.PersistentPreRunE != nil {
				if err := parent.PersistentPreRunE(parent, args); err != nil {
					return err
				}
			}
		}

		// this gets called twice due to the way Persistent* funcs are inheritted
		return viper.ReadInConfig()
	},
}

func GetMainCommand() *cobra.Command {
	return mainCmd
}
