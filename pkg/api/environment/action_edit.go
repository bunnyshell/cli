package environment

import (
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
	"github.com/spf13/pflag"
)

type EditOptions struct {
	DeployOptions

	sdk.EnvironmentEditAction

	EditData

	WithDeploy bool
}

type EditData struct {
	Name string

	RemoteDevelopmentAllowed bool
	AutoUpdate               bool

	K8SIntegration string
}

func NewEditOptions() *EditOptions {
	environmentEditAction := sdk.NewEnvironmentEditAction()

	return &EditOptions{
		DeployOptions: *NewDeployOptions(""),

		EditData: EditData{
			RemoteDevelopmentAllowed: true,
			AutoUpdate:               true,
		},

		EnvironmentEditAction: *environmentEditAction,
	}
}

func (eo *EditOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	data := eo.EditData

	flags.StringVar(&data.Name, "name", data.Name, "Update environment name")
	flags.BoolVar(&data.RemoteDevelopmentAllowed, "remote-development", data.RemoteDevelopmentAllowed, "Allow remote development for the environment")
	flags.BoolVar(&data.AutoUpdate, "auto-update", data.AutoUpdate, "Automatically update the environment when components git refs change")

	flags.BoolVar(&eo.WithDeploy, "deploy", eo.WithDeploy, "Deploy the environment after update")
	flags.StringVar(&data.K8SIntegration, "k8s", data.K8SIntegration, "Set Kubernetes integration for the environment (if not set)")

	eo.DeployOptions.UpdateFlagSet(flags)
}

func Edit(options *EditOptions) (*sdk.EnvironmentItem, error) {
	model, resp, err := EditRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, nil
}

func EditRaw(options *EditOptions) (*sdk.EnvironmentItem, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).EnvironmentApi.EnvironmentEdit(ctx, options.ID)

	return applyEditOptions(request, options).Execute()
}

func applyEditOptions(request sdk.ApiEnvironmentEditRequest, options *EditOptions) sdk.ApiEnvironmentEditRequest {
	if options.EditData.Name != "" {
		options.EnvironmentEditAction.SetName(options.EditData.Name)
	}

	if options.EditData.K8SIntegration != "" {
		options.EnvironmentEditAction.SetKubernetesIntegration(options.EditData.K8SIntegration)
	}

	options.EnvironmentEditAction.SetAutoUpdate(options.EditData.AutoUpdate)
	options.EnvironmentEditAction.SetRemoteDevelopmentAllowed(options.EditData.RemoteDevelopmentAllowed)

	request = request.EnvironmentEditAction(options.EnvironmentEditAction)

	return request
}
