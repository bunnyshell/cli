package remote_development

import (
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/k8s/bridge"
	"bunnyshell.com/cli/pkg/lib"
	remoteDevPkg "bunnyshell.com/cli/pkg/remote_development"
	"bunnyshell.com/cli/pkg/remote_development/action"
	upAction "bunnyshell.com/cli/pkg/remote_development/action/up"
	remoteDevConfig "bunnyshell.com/cli/pkg/remote_development/config"
	"github.com/spf13/cobra"
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	resourceLoader := bridge.NewResourceLoader()
	upOptions := upAction.NewOptions(remoteDevConfig.NewManager(), resourceLoader)

	noTTY := false

	command := &cobra.Command{
		Use: "up",

		ValidArgsFunction: cobra.NoFileCompletions,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := upOptions.Validate(); err != nil {
				return err
			}

			return lib.OnlyStylish(cmd, args)
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			upOptions.SetCommand(args)

			if err := resourceLoader.Load(settings.Profile); err != nil {
				return err
			}

			upParameters, err := upOptions.ToParameters()
			if err != nil {
				return err
			}

			upAction := action.NewUp(*resourceLoader.Environment)

			if err = upAction.Run(upParameters); err != nil {
				return err
			}

			sshConfigFile, _ := remoteDevPkg.GetSSHConfigFilePath()
			cmd.Println("Pod is ready for Remote Development.")
			cmd.Printf("You can find the SSH Config file in %s\n", sshConfigFile)

			// start
			if !noTTY {
				if err = upAction.StartSSHTerminal(); err != nil {
					return err
				}
			}

			return upAction.Wait()
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.Organization.GetFlag("organization"))
	flags.AddFlag(options.Project.GetFlag("project"))
	flags.AddFlag(options.Environment.GetFlag("environment"))
	flags.AddFlag(options.ServiceComponent.GetFlag("component"))

	upOptions.UpdateFlagSet(command, flags)

	flags.BoolVar(&noTTY, "no-tty", false, "Start remote development with no SSH terminal")

	mainCmd.AddCommand(command)
}
