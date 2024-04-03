package action

import (
	"errors"

	"bunnyshell.com/cli/pkg/api/environment"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/util"
	"github.com/spf13/cobra"
)

var (
	errK8SIntegrationNotProvided = errors.New("kubernetes integration must be provided when deploying")
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	createOptions := environment.NewCreateOptions()

	command := &cobra.Command{
		Use: "create",

		ValidArgsFunction: cobra.NoFileCompletions,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := createOptions.Validate(); err != nil {
				return err
			}

			if createOptions.WithDeploy && createOptions.GetKubernetesIntegration() == "" {
				if !settings.IsStylish() {
					return errK8SIntegrationNotProvided
				}
			}

			return validateActionOptions(&createOptions.ActionOptions)
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			createOptions.Project = settings.Profile.Context.Project

			if err := createOptions.AttachGenesis(); err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			model, err := environment.Create(createOptions)
			if err != nil {
				return createOptions.HandleError(cmd, err)
			}

			if !createOptions.WithDeploy {
				return lib.FormatCommandData(cmd, model)
			}

			deployOptions := &createOptions.DeployOptions
			deployOptions.ID = model.GetId()

			return HandleDeploy(cmd, deployOptions, "created", createOptions.GetKubernetesIntegration(), settings.IsStylish())
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.Project.AddFlagWithExtraHelp(
		"project",
		"Project for the environment",
		"Projects contain environments along with build settings and project variables",
		util.FlagRequired,
	))

	createOptions.UpdateCommandFlags(command)

	mainCmd.AddCommand(command)
}
