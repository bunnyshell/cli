package remote_development

import (
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/k8s/bridge"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/remote_development/action/up"
	remoteDevConfig "bunnyshell.com/cli/pkg/remote_development/config"
	"github.com/spf13/cobra"
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	resourceLoader := bridge.NewResourceLoader()
	upOptions := up.NewOptions(remoteDevConfig.NewManager(), resourceLoader)

	command := &cobra.Command{
		Use: "config",

		ValidArgsFunction: cobra.NoFileCompletions,

		PreRunE: lib.OnlyStylish,

		Hidden: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			upOptions.SetCommand(args)

			if err := resourceLoader.Load(settings.Profile); err != nil {
				return err
			}

			upParameters, err := upOptions.ToParameters()
			if err != nil {
				return err
			}

			return lib.FormatCommandData(cmd, upParameters)
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.Organization.GetFlag("organization"))
	flags.AddFlag(options.Project.GetFlag("project"))
	flags.AddFlag(options.Environment.GetFlag("environment"))
	flags.AddFlag(options.ServiceComponent.GetFlag("component"))

	upOptions.UpdateFlagSet(command, flags)

	mainCmd.AddCommand(command)
}
