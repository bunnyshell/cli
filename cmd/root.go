package cmd

import (
	"errors"
	"fmt"
	"os"

	"bunnyshell.com/cli/cmd/component"
	"bunnyshell.com/cli/cmd/configure"
	"bunnyshell.com/cli/cmd/environment"
	"bunnyshell.com/cli/cmd/event"
	"bunnyshell.com/cli/cmd/organization"
	"bunnyshell.com/cli/cmd/port_forward"
	"bunnyshell.com/cli/cmd/project"
	"bunnyshell.com/cli/cmd/remote_development"
	"bunnyshell.com/cli/cmd/variable"
	"bunnyshell.com/cli/cmd/version"
	"bunnyshell.com/cli/pkg/build"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/net"
	"bunnyshell.com/cli/pkg/util"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:     build.Name,
	Version: build.Version,

	Short: "Bunnyshell CLI",
	Long:  "Bunnyshell CLI helps you manage environments in Bunnyshell and enable Remote Development.",

	SilenceUsage: true,

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		manager := config.MainManager

		if cmd.CalledAs() == cobra.ShellCompRequestCmd {
			// Autocomplete parses flags differently, kickstart flag parsing
			_ = cmd.Root().ParseFlags(args)
			manager.Load()

			return nil
		}

		manager.Load()

		if errors.Is(manager.Error, config.ErrUnknownProfile) {
			return manager.Error
		}

		settings := config.GetSettings()

		if settings.NoProgress {
			net.DefaultSpinnerTransport.Disabled = true
		}
		if settings.Verbosity != 0 {
			fmt.Fprintf(os.Stdout, "Using config file: %s\n", config.GetSettings().ConfigFile)
		}

		cmd.SetOut(os.Stdout)
		cmd.SetErr(os.Stdout)

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	util.AddGroupedCommands(
		rootCmd,
		cobra.Group{
			ID:    "resources",
			Title: "Bunnyshell Resources",
		},
		[]*cobra.Command{
			component.GetMainCommand(),
			environment.GetMainCommand(),
			event.GetMainCommand(),
			organization.GetMainCommand(),
			port_forward.GetMainCommand(),
			project.GetMainCommand(),
			remote_development.GetMainCommand(),
			variable.GetMainCommand(),
		},
	)

	util.AddGroupedCommands(
		rootCmd,
		cobra.Group{
			ID:    "cli",
			Title: "CLI",
		},
		[]*cobra.Command{
			configure.GetMainCommand(),
			version.GetMainCommand(),
		},
	)
	rootCmd.SetHelpCommandGroupID("cli")
	rootCmd.SetCompletionCommandGroupID("cli")

	config.MainManager.CommandWithGlobalOptions(rootCmd)
	util.AllComandsHelpFlag(rootCmd)
}
