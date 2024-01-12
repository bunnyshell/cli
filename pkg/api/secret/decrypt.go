package secret

import (
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/api/common"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
	"github.com/spf13/pflag"
)

type DecryptOptions struct {
	common.Options

	sdk.SecretDecryptAction
}

func NewDecryptOptions() *DecryptOptions {
	return &DecryptOptions{}
}

func (options *DecryptOptions) UpdateSelfFlags(flags *pflag.FlagSet) {
	flags.StringVar(&options.Organization, "organization", options.Organization, "The Organization this secret was encrypted for")
}

func Decrypt(options *DecryptOptions) (*sdk.SecretDecryptedItem, error) {
	model, resp, err := DecryptRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, nil
}

func DecryptRaw(options *DecryptOptions) (*sdk.SecretDecryptedItem, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).
		SecretAPI.SecretDecrypt(ctx).
		SecretDecryptAction(options.SecretDecryptAction)

	return request.Execute()
}
