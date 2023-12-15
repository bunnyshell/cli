package action

import (
	"context"
	"errors"
	"fmt"
	"time"

	"bunnyshell.com/cli/pkg/api/project"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/util"
	"bunnyshell.com/sdk"
	"github.com/avast/retry-go/v4"
	"github.com/spf13/cobra"
)

const (
	BuildSettingStatusSuccess    string = "success"
	BuildSettingStatusValidating string = "validating"
	BuildSettingStatusError      string = "error"
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	editBuildSettingsOptions := project.NewEditBuildSettingsOptions("")

	command := &cobra.Command{
		Use: "update-build-settings",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			editBuildSettingsOptions.ID = settings.Profile.Context.Project

			_, err := project.EditBuildSettings(editBuildSettingsOptions)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			model, err := validateBuildSettings(editBuildSettingsOptions, settings.IsStylish())
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			return lib.FormatCommandData(cmd, model)
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.Project.GetFlag("id", util.FlagRequired))

	editBuildSettingsOptions.UpdateFlagSet(flags)

	mainCmd.AddCommand(command)
}

func validateBuildSettings(options *project.EditBuildSettingsOptions, showSpinner bool) (*sdk.ProjectItem, error) {
	if showSpinner {
		util.MakeSpinner("Validating the build settings...")
	}

	itemOptions := project.NewItemOptions(options.ID)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(options.ValidationTimeout)*time.Second)
	defer cancel()

	model, err := retry.DoWithData(
		func() (*sdk.ProjectItem, error) {
			model, err := project.Get(itemOptions)
			if err != nil {
				return nil, err
			}

			buildSettings := model.GetBuildSettings()

			// Return the model if the build settings are not set at all or they are already validated
			if !model.BuildSettings.IsSet() || buildSettings.GetLastStatus() == BuildSettingStatusSuccess {
				return model, nil
			}

			if buildSettings.GetLastStatus() == BuildSettingStatusError {
				return nil, retry.Unrecoverable(fmt.Errorf("the new build settings validation failed: %s", buildSettings.GetLastError()))
			}

			return nil, errors.New("the new build settings validation is in progress: %s")
		},
		retry.Context(ctx),
		retry.Attempts(0),
		retry.Delay(5*time.Second),
	)

	if err == context.DeadlineExceeded {
		return nil, fmt.Errorf("the new build settings validation timed out after %d seconds", options.ValidationTimeout)
	}

	return model, err
}
