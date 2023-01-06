package profile

import (
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/util"
	"github.com/spf13/cobra"
)

func init() {
	command := &cobra.Command{
		Use: "list",

		ValidArgsFunction: cobra.NoFileCompletions,

		PersistentPreRunE: util.PersistentPreRunChain,

		RunE: func(cmd *cobra.Command, args []string) error {
			result := map[string]interface{}{
				"file": config.GetSettings().ConfigFile,
			}

			if config.MainManager.Error != nil {
				result["error"] = config.MainManager.Error.Error()
			} else {
				result["data"] = config.GetConfig().Profiles
			}

			return lib.FormatCommandData(cmd, result)
		},
	}

	mainCmd.AddCommand(command)
}
