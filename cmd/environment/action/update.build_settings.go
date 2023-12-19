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
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	useProjectSettings := enum.BoolFalse

	editBuildSettingsOptions := environment.NewEditBuildSettingsOptions("")

	command := &cobra.Command{
		Use: "update-build-settings",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			editBuildSettingsOptions.ID = settings.Profile.Context.Environment

			if useProjectSettings == enum.BoolTrue {
				editBuildSettingsOptions.EditData.UseManagedCluster = enum.BoolFalse
				editBuildSettingsOptions.EditData.RegistryIntegration = ""
				editBuildSettingsOptions.Cpu = sdk.NullableString{}
				editBuildSettingsOptions.Memory = sdk.NullableInt32{}
			}

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

	useProjectSettingsFlag := enum.BoolFlag(
		&useProjectSettings,
		"use-project-settings",
		"Use the project build settings",
	)
	flags.AddFlag(useProjectSettingsFlag)
	useProjectSettingsFlag.NoOptDefVal = "true"

	editBuildSettingsOptions.UpdateFlagSet(flags)

	// use-project-settings excludes the other build settings flags for the cluster
	command.MarkFlagsMutuallyExclusive("use-project-settings", "use-managed-k8s")
	command.MarkFlagsMutuallyExclusive("use-project-settings", "k8s")
	command.MarkFlagsMutuallyExclusive("use-project-settings", "cpu")
	command.MarkFlagsMutuallyExclusive("use-project-settings", "memory")

	mainCmd.AddCommand(command)
}
