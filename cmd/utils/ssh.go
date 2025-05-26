package utils

import (
	"os"

	"bunnyshell.com/cli/cmd/component/action"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/util"
	"github.com/spf13/cobra"
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	sshOptions := action.SSHOptions{
		Shell: "/bin/sh",
	}

	command := &cobra.Command{
		Use: "ssh",

		Short: "SSH into a running container for a component",

		ValidArgsFunction: cobra.NoFileCompletions,

		Run: func(cmd *cobra.Command, args []string) {
			proxyArgs := []string{
				"components", "ssh",
				"--id", settings.Profile.Context.ServiceComponent,
			}

			if sshOptions.PodName != "" {
				proxyArgs = append(proxyArgs, "--pod", sshOptions.PodName)
			}

			if sshOptions.Container != "" {
				proxyArgs = append(proxyArgs, "--container", sshOptions.Container)
			}

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
				os.Exit(1)
			}
		},
	}

	flags := command.Flags()

	flags.AddFlag(
		options.ServiceComponent.AddFlag("component", "Service Component", util.FlagRequired),
	)

	sshOptions.UpdateFlagSet(flags)

	mainCmd.AddCommand(command)
}
