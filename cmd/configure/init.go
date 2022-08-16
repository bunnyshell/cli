package configure

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/util"
)

func init() {
	initConfigCommand := &cobra.Command{
		Use:   "init",
		Short: "Create a configuration file",
		RunE: func(cmd *cobra.Command, args []string) error {
			configFile := viper.ConfigFileUsed()

			if err := assertNoFile(configFile); err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			configFile = maybeChangeFile(configFile)
			if err := os.MkdirAll(filepath.Dir(configFile), os.FileMode(int(0700))); err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			viper.SetConfigPermissions(os.FileMode(int(0600)))
			viper.SetConfigFile(configFile)
			viper.Set("profiles", []string{})

			if err := viper.WriteConfig(); err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			if lib.CLIContext.Verbosity == 0 {
				return nil
			}

			return lib.FormatCommandData(cmd, map[string]interface{}{
				"message": "Config file created",
				"data":    viper.GetViper().ConfigFileUsed(),
			})
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			ok, err := util.Confirm("Continue with profile creation")
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			if !ok {
				return nil
			}

			root := mainCmd.Root()
			root.SetArgs([]string{"configure", "profiles", "add"})
			return root.Execute()
		},
	}

	mainCmd.AddCommand(initConfigCommand)
}

func maybeChangeFile(path string) string {
	return util.AskDefault("Choose file:", path, util.ConfigFileValidation)
}

func assertNoFile(path string) error {
	exists, err := fileExists(path)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("file already exists")
	}

	return nil
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
