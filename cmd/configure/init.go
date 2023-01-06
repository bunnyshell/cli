package configure

import (
	"fmt"
	"path/filepath"

	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/interactive"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

func init() {
	settings := config.GetSettings()

	initConfigCommand := &cobra.Command{
		Use:   "init",
		Short: "Create a configuration file",

		ValidArgsFunction: cobra.NoFileCompletions,

		Deprecated: "All configure commands will create a config file or overwrite an existing file.",

		RunE: func(cmd *cobra.Command, args []string) error {
			configFile, err := askForConfigFile(settings)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			settings.ConfigFile = configFile

			if err = config.MainManager.SafeSave(); err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			return lib.FormatCommandData(cmd, map[string]interface{}{
				"message": "Config file created",
				"data":    settings.ConfigFile,
			})
		},

		PostRunE: func(cmd *cobra.Command, args []string) error {
			if settings.NonInteractive {
				return nil
			}

			ok, err := interactive.Confirm("Continue with profile creation")
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			if !ok {
				return nil
			}

			root := cmd.Root()
			root.SetArgs([]string{"configure", "profiles", "add"})

			return root.Execute()
		},
	}

	mainCmd.AddCommand(initConfigCommand)
}

func askForConfigFile(settings *config.Settings) (string, error) {
	if settings.NonInteractive {
		return settings.ConfigFile, nil
	}

	return interactive.AskPath("Choose file:", settings.ConfigFile, requiredExtension(".json", ".yaml"))
}

func requiredExtension(extensions ...string) survey.Validator {
	return func(input interface{}) error {
		str, ok := input.(string)
		if !ok {
			return interactive.ErrInvalidValue
		}

		ext := filepath.Ext(str)
		for _, allowed := range extensions {
			if ext == allowed {
				return nil
			}
		}

		return fmt.Errorf("%w: extensions must be one of %v", interactive.ErrInvalidValue, extensions)
	}
}
