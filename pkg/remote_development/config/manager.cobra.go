package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type ShellCompletion func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective)

func (manager *Manager) UpdateFlagSet(command *cobra.Command, flags *pflag.FlagSet) {
	configFileParam := "rdev-configFile"
	configDirParam := "rdev-configDir"
	profileNameParam := "rdev-profile"

	flags.StringVar(
		&manager.configFileParam,
		configFileParam,
		manager.configFileParam,
		fmt.Sprintf(
			"Remote Dev config file\n"+
				"An absolute path ignores --%s",
			configDirParam,
		),
	)

	flags.StringVar(
		&manager.configDirParam,
		configDirParam,
		manager.configDirParam,
		"Remote Dev config directory\n"+
			"Using ... will look through all parent directories",
	)

	flags.StringVar(
		&manager.profileName,
		profileNameParam,
		manager.profileName,
		fmt.Sprintf(
			"Remote Dev profile name. Loaded from --%s",
			configFileParam,
		),
	)

	_ = command.RegisterFlagCompletionFunc(profileNameParam, manager.profileNamesCompletion())
}

func (manager *Manager) profileNamesCompletion() ShellCompletion {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if manager.Load() != nil {
			return []string{}, cobra.ShellCompDirectiveDefault
		}

		return manager.config.profileNames(), cobra.ShellCompDirectiveDefault
	}
}
