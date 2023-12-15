package environment

import (
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/util"
	"bunnyshell.com/sdk"
	"github.com/spf13/pflag"
)

type CreateOptions struct {
	DeployOptions

	sdk.EnvironmentCreateAction

	WithDeploy bool

	Labels map[string]string
}

func NewCreateOptions() *CreateOptions {
	environmentCreateAction := sdk.NewEnvironmentCreateActionWithDefaults()
	environmentCreateAction.SetKubernetesIntegration("")
	environmentCreateAction.SetEphemeralKubernetesIntegration("")

	return &CreateOptions{
		DeployOptions: *NewDeployOptions(""),

		EnvironmentCreateAction: *environmentCreateAction,
	}
}

func (co *CreateOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	k8sIntegration := co.KubernetesIntegration.Get()

	flags.StringVar(&co.Name, "name", co.Name, "Unique name for the environment")
	flags.BoolVar(&co.WithDeploy, "deploy", co.WithDeploy, "Deploy the environment after creation")
	flags.StringVar(k8sIntegration, "k8s", *k8sIntegration, "Use a Kubernetes integration for the environment")

	util.MarkFlagRequiredWithHelp(flags.Lookup("name"), "A unique name within the project for the new environment")

	flags.StringToStringVar(&co.Labels, "label", co.Labels, "Set labels for the new environment (key=value)")

	co.DeployOptions.UpdateFlagSet(flags)
}

func Create(options *CreateOptions) (*sdk.EnvironmentItem, error) {
	model, resp, err := CreateRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, nil
}

func CreateRaw(options *CreateOptions) (*sdk.EnvironmentItem, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).EnvironmentAPI.EnvironmentCreate(ctx)

	if len(options.Labels) > 0 {
		options.EnvironmentCreateAction.SetLabels(options.Labels)
	}

	request = request.EnvironmentCreateAction(options.EnvironmentCreateAction)

	return request.Execute()
}
