package project

import (
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/api/common"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/util"
	"bunnyshell.com/sdk"
	"github.com/spf13/pflag"
)

type EditBuildSettingsOptions struct {
	common.ItemOptions

	sdk.ProjectEditBuildSettingsAction

	EditBuildSettingsData

	// Seconds to wait for the build settings to be validated
	ValidationTimeout int32
}

type EditBuildSettingsData struct {
	UseManagedRegistry  bool
	RegistryIntegration string

	UseManagedCluster   bool
	BuildK8sIntegration string
}

func NewEditBuildSettingsOptions(project string) *EditBuildSettingsOptions {
	return &EditBuildSettingsOptions{
		ItemOptions: *common.NewItemOptions(project),

		EditBuildSettingsData: EditBuildSettingsData{},

		ProjectEditBuildSettingsAction: *sdk.NewProjectEditBuildSettingsAction(),

		ValidationTimeout: 180,
	}
}

func (eso *EditBuildSettingsOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	data := &eso.EditBuildSettingsData

	flags.BoolVar(&data.UseManagedRegistry, "use-managed-registry", data.UseManagedRegistry, "Use the managed Container Registry for the built images")
	flags.StringVar(&data.RegistryIntegration, "registry", data.RegistryIntegration, "Set the Container Registry integration to push the built images")

	flags.BoolVar(&data.UseManagedCluster, "use-managed-k8s", data.UseManagedCluster, "Use the managed Kubernetes integration cluster for the image builds")
	flags.StringVar(&data.BuildK8sIntegration, "k8s", data.BuildK8sIntegration, "Set the Kubernetes integration cluster to be used for the image builds")

	flags.Int32Var(&eso.ValidationTimeout, "validation-timeout", eso.ValidationTimeout, "Seconds to wait for the build settings to be validated")
}

func EditBuildSettings(options *EditBuildSettingsOptions) (*sdk.ProjectItem, error) {
	model, resp, err := EditBuildSettingsRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, nil
}

func EditBuildSettingsRaw(options *EditBuildSettingsOptions) (*sdk.ProjectItem, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).ProjectAPI.ProjectEditBuildSettings(ctx, options.ID)

	return applyEditBuildSettingsOptions(request, options).Execute()
}

func applyEditBuildSettingsOptions(
	request sdk.ApiProjectEditBuildSettingsRequest,
	options *EditBuildSettingsOptions,
) sdk.ApiProjectEditBuildSettingsRequest {
	if util.IsFlagPassed("use-managed-registry") {
		options.ProjectEditBuildSettingsAction.SetUseManagedRegistry(options.EditBuildSettingsData.UseManagedRegistry)
	}

	if options.EditBuildSettingsData.RegistryIntegration != "" {
		options.ProjectEditBuildSettingsAction.SetRegistryIntegration(options.EditBuildSettingsData.RegistryIntegration)
	}

	if util.IsFlagPassed("use-managed-k8s") {
		options.ProjectEditBuildSettingsAction.SetUseManagedCluster(options.EditBuildSettingsData.UseManagedCluster)
	}

	if options.EditBuildSettingsData.BuildK8sIntegration != "" {
		options.ProjectEditBuildSettingsAction.SetKubernetesIntegration(options.EditBuildSettingsData.BuildK8sIntegration)
	}

	request = request.ProjectEditBuildSettingsAction(options.ProjectEditBuildSettingsAction)

	return request
}
