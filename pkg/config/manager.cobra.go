package config

import (
	"github.com/spf13/cobra"
)

type ShellCompletion func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective)

func (manager *Manager) CommandWithGlobalOptions(command *cobra.Command) {
	flags := command.PersistentFlags()

	configFileFlag := manager.options.ConfigFile.GetMainFlag()
	flags.AddFlag(configFileFlag)
	_ = flags.SetAnnotation(configFileFlag.Name, cobra.BashCompFilenameExt, []string{"yaml", "json"})

	flags.AddFlag(manager.options.Debug.GetMainFlag())
	flags.AddFlag(manager.options.NoProgress.GetMainFlag())
	flags.AddFlag(manager.options.NonInteractive.GetMainFlag())
	flags.AddFlag(manager.options.Verbosity.GetMainFlag())

	profileFlag := manager.options.ProfileName.GetMainFlag()
	flags.AddFlag(profileFlag)
	_ = command.RegisterFlagCompletionFunc(profileFlag.Name, manager.profileNamesCompletion())

	outputFormatFlag := manager.options.OutputFormat.GetMainFlag()
	flags.AddFlag(outputFormatFlag)
	_ = command.RegisterFlagCompletionFunc(outputFormatFlag.Name, manager.outputTypesCompletion())
}

func (manager *Manager) CommandWithAPI(command *cobra.Command) {
	flags := command.PersistentFlags()

	tokenFlag := manager.options.Token.GetMainFlag()
	flags.AddFlag(tokenFlag)
	_ = command.MarkPersistentFlagRequired(tokenFlag.Name)

	flags.AddFlag(manager.options.Host.GetMainFlag())
	flags.AddFlag(manager.options.Timeout.GetMainFlag())
}

func (manager *Manager) profileNamesCompletion() ShellCompletion {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		MainManager.Load()

		return MainManager.config.profileNames(), cobra.ShellCompDirectiveNoFileComp
	}
}

func (manager *Manager) outputTypesCompletion() ShellCompletion {
	return cobra.FixedCompletions(FormatDescriptions, cobra.ShellCompDirectiveDefault)
}
