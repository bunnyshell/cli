package action

import (
	"bunnyshell.com/cli/pkg/api/environment"
	"bunnyshell.com/cli/pkg/config"
	"github.com/spf13/cobra"
)

type DeployData struct {
	K8SIntegration string
}

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	deployOptions := environment.NewDeployOptions("")
	deployData := DeployData{}

	command := &cobra.Command{
		Use: "deploy",

		ValidArgsFunction: cobra.NoFileCompletions,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return validateActionOptions(&deployOptions.ActionOptions)
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			deployOptions.ID = settings.Profile.Context.Environment

			return handleDeploy(cmd, deployOptions, "", deployData.K8SIntegration)
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.Environment.GetRequiredFlag("id"))

	deployOptions.UpdateFlagSet(flags)

	flags.StringVar(&deployData.K8SIntegration, "k8s", deployData.K8SIntegration, "Use a Kubernetes integration for the deployment (if not set)")

	mainCmd.AddCommand(command)
}
