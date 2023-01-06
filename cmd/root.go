package cmd

import (
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
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/lib/cliconfig"
	"bunnyshell.com/cli/pkg/net"
	"bunnyshell.com/cli/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:     build.Name,
	Version: build.Version,

	Short: "Bunnyshell CLI",
	Long:  "Bunnyshell CLI helps you manage environments in Bunnyshell and enable Remote Development.",

	SilenceUsage: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cmd.CalledAs() == cobra.ShellCompRequestCmd {
			return
		}

		cmd.SetOut(os.Stdout)
		cmd.SetErr(os.Stdout)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

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

	lib.CLIContext.SetGlobalFlags(rootCmd)
}

func initConfig() {
	if lib.CLIContext.NoProgress {
		net.DefaultSpinnerTransport.Disabled = true
	}

	cobra.CheckErr(cliconfig.FindConfigFile())

	viper.SetEnvPrefix(build.EnvPrefix)
	viper.AutomaticEnv()

	if lib.CLIContext.Verbosity != 0 {
		fmt.Fprintln(os.Stdout, "Using config file:", viper.ConfigFileUsed())
	}
}
