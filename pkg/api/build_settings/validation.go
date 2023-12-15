package build_settings

import (
	"context"
	"errors"
	"fmt"
	"time"

	"bunnyshell.com/cli/pkg/api/common"
	"bunnyshell.com/cli/pkg/util"
	"bunnyshell.com/sdk"
	"github.com/avast/retry-go/v4"
)

type ModelWithBuildSettings[T any] interface {
	HasBuildSettings() bool
	GetBuildSettings() sdk.BuildSettingsItem
	*T
}

type ModelFetcher[T any, PT ModelWithBuildSettings[T]] func(*common.ItemOptions) (PT, error)

func CheckBuildSettingsValidation[T any, PT ModelWithBuildSettings[T]](fetcher ModelFetcher[T, PT], options *EditOptions, showSpinner bool) (PT, error) {
	if showSpinner {
		util.MakeSpinner("Validating the build settings...")
	}

	itemOptions := common.NewItemOptions(options.ID)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(options.ValidationTimeout)*time.Second)
	defer cancel()

	model, err := retry.DoWithData(
		func() (PT, error) {
			model, err := fetcher(itemOptions)
			if err != nil {
				return nil, err
			}

			buildSettings := model.GetBuildSettings()

			// Return the model if the build settings are not set at all or they are already validated
			if !model.HasBuildSettings() || buildSettings.GetLastStatus() == StatusSuccess {
				return model, nil
			}

			if buildSettings.GetLastStatus() == StatusError {
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
