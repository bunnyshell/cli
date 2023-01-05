package cliconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"bunnyshell.com/cli/pkg/build"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/util"
)

func FindConfigFile() error {
	configFile := determineConfig()

	configDotExtension := filepath.Ext(configFile)
	if configDotExtension == "" {
		return fmt.Errorf("unsuported extension")
	}
	configExtension := configDotExtension[1:]

	if !stringInSlice(configExtension, viper.SupportedExts) {
		return fmt.Errorf("unsuported extension %s", configExtension)
	}

	viper.SetConfigFile(configFile)

	return nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func determineConfig() string {
	if lib.CLIContext.ConfigFile != "" {
		return lib.CLIContext.ConfigFile
	}

	cfgFile, ok := os.LookupEnv(strings.ToUpper(build.EnvPrefix) + "_CONFIG_FILE")
	if ok {
		return cfgFile
	}

	workspace, err := util.GetWorkspaceDir()
	cobra.CheckErr(err)

	return workspace + "/config.yaml"
}
