package template

import (
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/api/common"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
	"github.com/spf13/pflag"
)

type ValidateOptions struct {
	common.Options

	sdk.TemplateValidateAction

	Organization string

	WithComponents bool
}

func (vo *ValidateOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	flags.BoolVar(&vo.WithComponents, "with-components", vo.WithComponents, "Validate components along with the template")
}

func NewValidateOptions() *ValidateOptions {
	return &ValidateOptions{
		Options: *common.NewOptions(),

		TemplateValidateAction: sdk.TemplateValidateAction{},

		WithComponents: false,
	}
}

func Validate(options *ValidateOptions) (bool, error) {
	_, resp, err := ValidateRaw(options)
	if err != nil {
		return false, api.ParseError(resp, err)
	}

	return true, nil
}

func ValidateRaw(options *ValidateOptions) (*sdk.TemplateCollection, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).TemplateApi.TemplateValidate(ctx)

	return applyValidateOptions(request, options).Execute()
}

func applyValidateOptions(request sdk.ApiTemplateValidateRequest, options *ValidateOptions) sdk.ApiTemplateValidateRequest {
	request = request.TemplateValidateAction(options.TemplateValidateAction)

	return request
}
