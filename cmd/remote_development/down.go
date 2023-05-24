package remote_development

import (
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/k8s/bridge"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/remote_development/action"
	"bunnyshell.com/cli/pkg/remote_development/action/down"
	remoteDevConfig "bunnyshell.com/cli/pkg/remote_development/config"
	"github.com/spf13/cobra"
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	resourceLoader := bridge.NewResourceLoader()
	downOptions := down.NewOptions(remoteDevConfig.NewManager(), resourceLoader)

	command := &cobra.Command{
		Use: "down",

		ValidArgsFunction: cobra.NoFileCompletions,

		PreRunE: lib.OnlyStylish,

		RunE: func(cmd *cobra.Command, args []string) error {
			if err := resourceLoader.Load(settings.Profile); err != nil {
				return err
			}

			downParameters, err := downOptions.ToParameters()
			if err != nil {
				return err
			}

			downAction := action.NewDown(*resourceLoader.Environment)

			return downAction.Run(downParameters)
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.Organization.GetFlag("organization"))
	flags.AddFlag(options.Project.GetFlag("project"))
	flags.AddFlag(options.Environment.GetFlag("environment"))
	flags.AddFlag(options.ServiceComponent.GetFlag("component"))

	downOptions.UpdateFlagSet(command, flags)

	mainCmd.AddCommand(command)
}
