package project

import (
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/api/common"
	"bunnyshell.com/cli/pkg/lib"
)

type DeleteOptions struct {
	common.ItemOptions
}

func NewDeleteOptions() *DeleteOptions {
	return &DeleteOptions{}
}

func Delete(options *DeleteOptions) error {
	resp, err := DeleteRaw(options)
	if err != nil {
		return api.ParseError(resp, err)
	}

	return nil
}

func DeleteRaw(options *DeleteOptions) (*http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).
		ProjectAPI.ProjectDelete(ctx, options.ID)

	return request.Execute()
}
