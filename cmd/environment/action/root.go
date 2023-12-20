package action

import (
	"errors"
	"fmt"

	"bunnyshell.com/cli/pkg/api/common"
	"bunnyshell.com/cli/pkg/api/component/endpoint"
	"bunnyshell.com/cli/pkg/api/environment"
	"bunnyshell.com/cli/pkg/build"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/interactive"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/progress"
	"bunnyshell.com/sdk"
	"github.com/spf13/cobra"
)

var mainCmd = &cobra.Command{}

var (
	errOtherK8s = errors.New("environment has a different Kubernetes Integration")

	errK8sRequired = errors.New("environment requires a Kubernetes Integration for deployment")
)

func GetMainCommand() *cobra.Command {
	return mainCmd
}

func validateActionOptions(actionOptions *common.ActionOptions) error {
	if !actionOptions.WithoutPipeline {
		return nil
	}

	if config.GetSettings().IsStylish() {
		return nil
	}

	return fmt.Errorf("%w when following pipelines", lib.ErrNotStylish)
}

func HandleDeploy(cmd *cobra.Command, deployOptions *environment.DeployOptions, action string, kubernetesIntegration string) error {
	if err := ensureKubernetesIntegration(deployOptions, kubernetesIntegration); err != nil {
		return err
	}

	if action != "" {
		cmd.Printf("\nEnvironment %s successfully %s... deploying...\n", deployOptions.ID, action)
	}

	event, err := environment.Deploy(deployOptions)
	if err != nil {
		return lib.FormatCommandError(cmd, err)
	}

	if deployOptions.WithoutPipeline {
		return lib.FormatCommandData(cmd, event)
	}

	if err = processEventPipeline(cmd, event, "deploy"); err != nil {
		cmd.Printf("\nEnvironment %s deploying failed\n", deployOptions.ID)

		return err
	}

	cmd.Printf("\nEnvironment %s successfully deployed\n", deployOptions.ID)

	return showEnvironmentEndpoints(cmd, deployOptions.ID)
}

func ensureKubernetesIntegration(deployOptions *environment.DeployOptions, kubernetesIntegration string) error {
	model, err := environment.Get(environment.NewItemOptions(deployOptions.ID))
	if err != nil {
		return err
	}

	environmentKubernetesIntegration := model.GetKubernetesIntegration()
	if environmentKubernetesIntegration != "" {
		if kubernetesIntegration == "" {
			return nil
		}

		if environmentKubernetesIntegration == kubernetesIntegration {
			return nil
		}

		return fmt.Errorf("%w: %s", errOtherK8s, environmentKubernetesIntegration)
	}

	if kubernetesIntegration == "" {
		if !config.GetSettings().IsStylish() {
			return fmt.Errorf("%w and cannot be interactively supplied in non-stylish mode", errK8sRequired)
		}

		question := interactive.NewInput("Deployment requires a Kubernetes Integration:")
		question.SetValidate(interactive.AssertMinimumLength(1))
		question.Help = fmt.Sprintf(
			`Find available Kubernetes Integrations with "%s k8s list"`,
			build.Name,
		)

		kubernetesIntegration, err = question.AskString()
		if err != nil {
			return err
		}
	}

	editSettingsOptions := environment.NewEditSettingsOptions(deployOptions.ID)
	editSettingsOptions.UpdateEditSettingsForType(model.GetType())

	editSettingsOptions.EnvironmentEditSettings.KubernetesIntegration.Set(&kubernetesIntegration)

	_, err = environment.EditSettings(editSettingsOptions)

	return err
}

func processEventPipeline(cmd *cobra.Command, event *sdk.EventItem, action string) error {
	progressOptions := progress.NewOptions()

	cmd.Printf(
		"Environment %s scheduled to %s with EventID %s\n",
		event.GetEnvironment(),
		action,
		event.GetId(),
	)

	pipeline, err := progress.EventToPipeline(event, progressOptions)
	if err != nil {
		return err
	}

	cmd.Printf(
		"EventID %s generated %s pipeline %s\n",
		pipeline.GetEvent(),
		action,
		pipeline.GetId(),
	)

	if err = progress.Pipeline(pipeline.GetId(), nil); err != nil {
		return err
	}

	return nil
}

func showEnvironmentEndpoints(cmd *cobra.Command, environment string) error {
	options := endpoint.NewAggregateOptions()
	options.Environment = environment

	components, err := endpoint.Aggregate(options)
	if err != nil {
		return err
	}

	return lib.FormatCommandData(cmd, components)
}
