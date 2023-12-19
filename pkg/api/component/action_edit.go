package component

import (
	"bunnyshell.com/cli/pkg/api/common"
	"bunnyshell.com/cli/pkg/api/environment"
	"github.com/spf13/pflag"
)

type EditComponentOptions struct {
	common.ItemOptions

	environment.DeployOptions

	EditComponentData

	WithDeploy bool
}

type EditComponentData struct {
	K8SIntegration string

	TargetRepository string
	TargetBranch     string
}

func NewEditComponentOptions(component string) *EditComponentOptions {
	return &EditComponentOptions{
		ItemOptions: *common.NewItemOptions(component),

		DeployOptions: *environment.NewDeployOptions(""),
	}
}

func (eo *EditComponentOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	data := eo.EditComponentData

	flags.BoolVar(&eo.WithDeploy, "deploy", eo.WithDeploy, "Deploy the environment after update")

	flags.StringVar(&data.K8SIntegration, "k8s", data.K8SIntegration, "Set Kubernetes integration for the environment (if not set)")

	eo.DeployOptions.UpdateFlagSet(flags)
}
