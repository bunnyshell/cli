package environment

import (
	"io"
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/api/common"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
)

type KubeConfigData map[string]any

type KubeConfigItem struct {
	Data *sdk.EnvironmentKubeConfigKubeConfigRead

	Bytes []byte
}

type KubeConfigOptions struct {
	common.ItemOptions
}

func NewKubeConfigOptions(id string) *KubeConfigOptions {
	return &KubeConfigOptions{
		ItemOptions: *common.NewItemOptions(id),
	}
}

func KubeConfig(options *KubeConfigOptions) (*KubeConfigItem, error) {
	model, resp, err := KubeConfigRaw(&options.ItemOptions)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &KubeConfigItem{
		Data:  model,
		Bytes: bytes,
	}, nil
}

func KubeConfigRaw(options *common.ItemOptions) (*sdk.EnvironmentKubeConfigKubeConfigRead, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).EnvironmentApi.EnvironmentKubeConfig(ctx, options.ID)

	return request.Execute()
}
