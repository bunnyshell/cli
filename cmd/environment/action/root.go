package action

import (
	"fmt"

	"bunnyshell.com/cli/pkg/api/common"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/progress"
	"bunnyshell.com/sdk"
	"github.com/spf13/cobra"
)

var mainCmd = &cobra.Command{}

func GetMainCommand() *cobra.Command {
	return mainCmd
}

func validateActionOptions(actionOptions *common.ActionOptions) error {
	if !actionOptions.WithPipeline {
		return nil
	}

	if config.GetSettings().IsStylish() {
		return nil
	}

	return fmt.Errorf("%w when following deployments", lib.ErrNotStylish)
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
		return lib.FormatCommandError(cmd, err)
	}

	cmd.Printf(
		"Event %s generated %s pipeline %s\n",
		pipeline.GetEvent(),
		action,
		pipeline.GetId(),
	)

	if err = progress.Pipeline(pipeline.GetId(), nil); err != nil {
		return lib.FormatCommandError(cmd, err)
	}

	return nil
}
