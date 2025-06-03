package environment

import (
	"github.com/spf13/pflag"
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/api/common"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
)

type DeleteOptions struct {
	common.ActionOptions

	QueueIfSomethingInProgress bool
}

func NewDeleteOptions(id string) *DeleteOptions {
	return &DeleteOptions{
		ActionOptions:              *common.NewActionOptions(id),
		QueueIfSomethingInProgress: false,
	}
}

func (options *DeleteOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	options.ActionOptions.UpdateFlagSet(flags)

	//flags.BoolVar(&options.QueueIfSomethingInProgress, "queue", options.QueueIfSomethingInProgress, "Queue the delete pipeline if another operation is in progress now")
}

func Delete(options *DeleteOptions) (*sdk.EventItem, error) {
	model, resp, err := DeleteRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, nil
}

func DeleteRaw(options *DeleteOptions) (*sdk.EventItem, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).EnvironmentAPI.EnvironmentDelete(ctx, options.ID)

	return request.Execute()
}
