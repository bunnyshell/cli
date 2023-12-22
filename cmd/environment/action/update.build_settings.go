package action

import (
	"bunnyshell.com/cli/pkg/api/build_settings"
	"bunnyshell.com/cli/pkg/api/environment"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/config/enum"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/util"
	"bunnyshell.com/sdk"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	useClusterProjectSettings := enum.BoolFalse
	useRegistryProjectSettings := enum.BoolFalse

	editBuildSettingsOptions := environment.NewEditBuildSettingsOptions("")

	command := &cobra.Command{
		Use: "update-build-settings",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			editBuildSettingsOptions.ID = settings.Profile.Context.Environment

			parseEditBuildSettingsOptions(cmd.Flags(), editBuildSettingsOptions, useClusterProjectSettings, useRegistryProjectSettings)

			_, err := environment.EditBuildSettings(editBuildSettingsOptions)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			model, err := build_settings.CheckBuildSettingsValidation[sdk.EnvironmentItem](
				environment.Get,
				&editBuildSettingsOptions.EditOptions,
				settings.IsStylish(),
			)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			return lib.FormatCommandData(cmd, model)
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.Environment.GetFlag("id", util.FlagRequired))

	useClusterProjectSettingsFlag := enum.BoolFlag(
		&useClusterProjectSettings,
		"use-project-k8s",
		"Use the project build cluster settings",
	)
	flags.AddFlag(useClusterProjectSettingsFlag)
	useClusterProjectSettingsFlag.NoOptDefVal = "true"

	useRegistryProjectSettingsFlag := enum.BoolFlag(
		&useRegistryProjectSettings,
		"use-project-registry",
		"Use the project build registry settings",
	)
	flags.AddFlag(useRegistryProjectSettingsFlag)
	useRegistryProjectSettingsFlag.NoOptDefVal = "true"

	editBuildSettingsOptions.UpdateFlagSet(flags)

	// use-project-settings excludes the other build settings flags for the cluster
	command.MarkFlagsMutuallyExclusive("use-project-k8s", "use-managed-k8s")
	command.MarkFlagsMutuallyExclusive("use-project-k8s", "k8s")
	command.MarkFlagsMutuallyExclusive("use-project-k8s", "cpu")
	command.MarkFlagsMutuallyExclusive("use-project-k8s", "memory")

	command.MarkFlagsMutuallyExclusive("use-project-registry", "use-managed-registry")
	command.MarkFlagsMutuallyExclusive("use-project-registry", "registry")

	mainCmd.AddCommand(command)
}

func parseEditBuildSettingsOptions(
	flags *pflag.FlagSet,
	editBuildSettingsOptions *environment.EditBuildSettingsOptions,
	useClusterProjectSettings enum.Bool,
	useRegistryProjectSettings enum.Bool,
) {
	if useClusterProjectSettings == enum.BoolTrue {
		editBuildSettingsOptions.EditData.UseManagedCluster = enum.BoolFalse
		editBuildSettingsOptions.SetKubernetesIntegration("")
		editBuildSettingsOptions.Cpu = sdk.NullableString{}
		editBuildSettingsOptions.Memory = sdk.NullableInt32{}
	}

	if useRegistryProjectSettings == enum.BoolTrue {
		editBuildSettingsOptions.EditData.UseManagedRegistry = enum.BoolFalse
		editBuildSettingsOptions.SetRegistryIntegration("")
	}
}
