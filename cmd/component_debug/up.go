package component_debug

import (
	"fmt"

	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/debug_component/action"
	upAction "bunnyshell.com/cli/pkg/debug_component/action/up"
	"bunnyshell.com/cli/pkg/k8s/bridge"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	sshOptions := SSHOptions{
		Shell:                 "/bin/sh",
		OverrideClusterServer: "",
	}

	resourceLoader := bridge.NewResourceLoader()
	upOptions := upAction.NewOptions(resourceLoader)

	command := &cobra.Command{
		Use: "start",

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

			selectedContainerName, err := upAction.GetSelectedContainerName()
			if err != nil {
				return err
			}

			sshOptions.OverrideClusterServer = upParameters.OverrideClusterServer
			if err = startSSH(*resourceLoader.Component.Id, selectedContainerName, sshOptions, cmd, args); err != nil {
				return fmt.Errorf("debug SSH exited with: %s", err)
			}

			return upAction.Close()
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.Organization.GetFlag("organization"))
	flags.AddFlag(options.Project.GetFlag("project"))
	flags.AddFlag(options.Environment.GetFlag("environment"))
	flags.AddFlag(options.ServiceComponent.GetFlag("component"))

	upOptions.UpdateFlagSet(command, flags)

	sshOptions.UpdateFlagSet(flags)

	mainCmd.AddCommand(command)
}

func startSSH(componentId string, containerName string, sshOptions SSHOptions, cmd *cobra.Command, args []string) error {
	proxyArgs := []string{
		"components", "ssh",
		"--id", componentId,
	}

	proxyArgs = append(proxyArgs, "--container", containerName)

	if sshOptions.Shell != "" {
		proxyArgs = append(proxyArgs, "--shell", sshOptions.Shell)
	}

	if sshOptions.NoBanner {
		proxyArgs = append(proxyArgs, "--no-banner")
	}

	if sshOptions.NoTTY {
		proxyArgs = append(proxyArgs, "--no-tty")
	}

	if sshOptions.OverrideClusterServer != "" {
		proxyArgs = append(proxyArgs, "--override-kubeconfig-cluster-server", sshOptions.OverrideClusterServer)
	}

	root := cmd.Root()
	root.SetArgs(append(proxyArgs, args...))

	if err := root.Execute(); err != nil {
		return err
	}

	return nil
}
