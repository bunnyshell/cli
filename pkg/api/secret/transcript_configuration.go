package secret

import (
	"io"
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/api/common"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
	"github.com/spf13/pflag"
)

type DefinitionData map[string]any

type DefinitionItem struct {
	Data DefinitionData

	Bytes []byte
}

type TranscriptMode string

const (
	TranscriptModeExposed    TranscriptMode = "exposed"
	TranscriptModeObfuscated TranscriptMode = "obfuscated"
	TranscriptModeResolved   TranscriptMode = "resolved"
)

type TranscriptConfigurationOptions struct {
	common.Options

	DefinitionFilePath string

	sdk.SecretTranscryptConfigurationAction
}

func NewTranscriptConfigurationOptions() *TranscriptConfigurationOptions {
	return &TranscriptConfigurationOptions{}
}

func (options *TranscriptConfigurationOptions) UpdateSelfFlags(flags *pflag.FlagSet) {
	flags.StringVar(&options.Organization, "organization", options.Organization, "The Organization this secret was encrypted for")
	flags.StringVar(&options.DefinitionFilePath, "file", options.DefinitionFilePath, "The filepath to the environment definition file")
}

func TranscriptConfiguration(options *TranscriptConfigurationOptions) (*DefinitionItem, error) {
	model, resp, err := TranscriptConfigurationRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &DefinitionItem{
		Data:  model,
		Bytes: bytes,
	}, nil
}

func TranscriptConfigurationRaw(options *TranscriptConfigurationOptions) (DefinitionData, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).
		SecretAPI.SecretTranscryptConfiguration(ctx).
		SecretTranscryptConfigurationAction(options.SecretTranscryptConfigurationAction)

	return request.Execute()
}
