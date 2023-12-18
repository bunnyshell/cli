package build_settings

import (
	"bunnyshell.com/cli/pkg/api/common"
	"bunnyshell.com/cli/pkg/config/enum"
	"github.com/spf13/pflag"
)

const (
	StatusSuccess    string = "success"
	StatusValidating string = "validating"
	StatusError      string = "error"
)

type ActionWithBuildSettings interface {
	SetUseManagedRegistry(bool)

	SetRegistryIntegration(string)

	SetUseManagedCluster(bool)

	SetKubernetesIntegration(string)
}

type EditOptions struct {
	common.ItemOptions

	EditData

	// Seconds to wait for the build settings to be validated
	ValidationTimeout int32
}

type EditData struct {
	UseManagedRegistry  enum.Bool
	RegistryIntegration string

	UseManagedCluster   enum.Bool
	BuildK8sIntegration string
}

func NewEditOptions(entityId string) *EditOptions {
	return &EditOptions{
		ItemOptions: *common.NewItemOptions(entityId),

		EditData: EditData{},

		ValidationTimeout: 180,
	}
}

func (eso *EditOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	data := &eso.EditData

	useManagedRegistryFlag := enum.BoolFlag(
		&data.UseManagedRegistry,
		"use-managed-registry",
		"Use the managed Container Registry for the built images",
	)
	flags.AddFlag(useManagedRegistryFlag)
	useManagedRegistryFlag.NoOptDefVal = "true"

	flags.StringVar(&data.RegistryIntegration, "registry", data.RegistryIntegration, "Set the Container Registry integration to push the built images")

	useManagedClusterFlag := enum.BoolFlag(
		&data.UseManagedCluster,
		"use-managed-k8s",
		"Use the managed Kubernetes integration cluster for the image builds",
	)
	flags.AddFlag(useManagedClusterFlag)
	useManagedClusterFlag.NoOptDefVal = "true"

	flags.StringVar(&data.BuildK8sIntegration, "k8s", data.BuildK8sIntegration, "Set the Kubernetes integration cluster to be used for the image builds")

	flags.Int32Var(&eso.ValidationTimeout, "validation-timeout", eso.ValidationTimeout, "Seconds to wait for the build settings to be validated")
}

func ApplyEditOptionsToAction(action ActionWithBuildSettings, options *EditData) {
	if options.UseManagedRegistry != enum.BoolNone {
		action.SetUseManagedRegistry(options.UseManagedRegistry == enum.BoolTrue)
	}

	if options.RegistryIntegration != "" {
		action.SetRegistryIntegration(options.RegistryIntegration)
	}

	if options.UseManagedCluster != enum.BoolNone {
		action.SetUseManagedCluster(options.UseManagedCluster == enum.BoolTrue)
	}

	if options.BuildK8sIntegration != "" {
		action.SetKubernetesIntegration(options.BuildK8sIntegration)
	}
}
