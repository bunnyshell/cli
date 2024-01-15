package secret

import (
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/api/common"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
	"github.com/spf13/pflag"
)

type EncryptOptions struct {
	common.Options

	sdk.SecretEncryptAction
}

func NewEncryptOptions() *EncryptOptions {
	return &EncryptOptions{}
}

func (options *EncryptOptions) UpdateSelfFlags(flags *pflag.FlagSet) {
	flags.StringVar(&options.Organization, "organization", options.Organization, "The Organization for which to encrypt the secret")
}

func Encrypt(options *EncryptOptions) (*sdk.SecretEncryptedItem, error) {
	model, resp, err := EncryptRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, nil
}

func EncryptRaw(options *EncryptOptions) (*sdk.SecretEncryptedItem, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).
		SecretAPI.SecretEncrypt(ctx).
		SecretEncryptAction(options.SecretEncryptAction)

	return request.Execute()
}
