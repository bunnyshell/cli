package environment

import (
	"io"
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/api/common"
	"bunnyshell.com/cli/pkg/lib"
)

type DefinitionData map[string]any

type DefinitionItem struct {
	Data DefinitionData

	Bytes []byte
}

type DefinitionOptions struct {
	common.ActionOptions
}

func NewDefinitionOptions(id string) *DefinitionOptions {
	return &DefinitionOptions{
		ActionOptions: *common.NewActionOptions(id),
	}
}

func Definition(options *DefinitionOptions) (*DefinitionItem, error) {
	model, resp, err := DefinitionRaw(options)
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

func DefinitionRaw(options *DefinitionOptions) (DefinitionData, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).EnvironmentApi.EnvironmentDefinition(ctx, options.ID)

	return request.Execute()
}
