package environment

import (
	"net/http"

	"github.com/spf13/pflag"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/api/common"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
)

type StopOptions struct {
	common.PartialActionOptions

	QueueIfSomethingInProgress bool
}

func NewStopOptions(id string) *StopOptions {
	return &StopOptions{
		PartialActionOptions:       *common.NewPartialActionOptions(id),
		QueueIfSomethingInProgress: false,
	}
}

func (options *StopOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	options.PartialActionOptions.UpdateFlagSet(flags)

	flags.BoolVar(&options.QueueIfSomethingInProgress, "queue", options.QueueIfSomethingInProgress, "Queue the stop pipeline if another operation is in progress now")
}

func Stop(options *StopOptions) (*sdk.EventItem, error) {
	model, resp, err := StopRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, nil
}

func StopRaw(options *StopOptions) (*sdk.EventItem, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	isPartialAction := options.IsPartial()

	request := lib.GetAPIFromProfile(profile).EnvironmentAPI.EnvironmentStop(ctx, options.ID).
		EnvironmentPartialStopAction(sdk.EnvironmentPartialStopAction{
			IsPartial:                  &isPartialAction,
			Components:                 options.GetActionComponents(),
			QueueIfSomethingInProgress: &options.QueueIfSomethingInProgress,
		})

	return request.Execute()
}
