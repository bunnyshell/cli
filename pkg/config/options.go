package config

import (
	"fmt"
	"strings"

	"bunnyshell.com/cli/pkg/build"
	"bunnyshell.com/cli/pkg/config/option"
	"bunnyshell.com/cli/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Options struct {
	// other options
	ConfigFile     *option.String
	Verbosity      *option.Count
	Timeout        *option.Duration
	NoProgress     *option.Bool
	NonInteractive *option.Bool

	// global options
	Debug        *option.Bool
	OutputFormat *option.String
	ProfileName  *option.String

	// profile options
	Host  *option.String
	Token *option.String

	// Profile.Context options
	Organization     *option.String
	Project          *option.String
	Environment      *option.String
	ServiceComponent *option.String
}

func NewOptions(settings *Settings) *Options {
	return &Options{
		ConfigFile: newConfigFile(settings),

		Verbosity:      newVerbosity(settings),
		Timeout:        newTimeout(settings),
		NoProgress:     newNoProgress(settings),
		NonInteractive: newNonInteractive(settings),

		Debug:        newDebug(settings),
		OutputFormat: newOutputFormat(settings),
		ProfileName:  newProfileName(settings),

		Token: newToken(settings),
		Host:  newHost(settings),

		Organization:     newOrganization(settings),
		Project:          newProject(settings),
		Environment:      newEnvironment(settings),
		ServiceComponent: newServiceComponent(settings),
	}
}

func newConfigFile(settings *Settings) *option.String {
	option := option.NewStringOption(&settings.ConfigFile)

	flag := option.AddFlag("configFile", "Bunnyshell CLI Config File")
	flag.Annotations = map[string][]string{
		cobra.BashCompFilenameExt: {"yaml", "json"},
	}

	if workspace, short, err := util.GetWorkspaceDirAndShort(); err == nil {
		_ = flag.Value.Set(workspace + "/config.yaml")
		flag.DefValue = short + "/config.yaml"
	}

	return option
}

func newVerbosity(settings *Settings) *option.Count {
	option := option.NewCountOption(&settings.Verbosity)

	option.AddFlagShort("verbose", "v", "Increase log verbosity")

	return option
}

func newTimeout(settings *Settings) *option.Duration {
	option := option.NewDurationOption(&settings.Timeout)

	option.AddFlagShort("timeout", "t", "Timeout value for network requests")

	return option
}

func newDebug(settings *Settings) *option.Bool {
	option := option.NewBoolOption(&settings.Debug)

	option.AddFlagShort("debug", "d", "Debug network requests")

	return option
}

func newOutputFormat(settings *Settings) *option.String {
	option := option.NewStringOption(&settings.OutputFormat)

	formatsString := strings.Join(Formats, " | ")

	option.Var().Validator = func(data string, flag pflag.Value) error {
		for _, format := range Formats {
			if format == data {
				return nil
			}
		}

		return fmt.Errorf("%w, expecting one of %s", ErrInvalidValue, formatsString)
	}

	option.AddFlagShort("output", "o", fmt.Sprintf("Output format: %s", formatsString))

	return option
}

func newNoProgress(settings *Settings) *option.Bool {
	option := option.NewBoolOption(&settings.NoProgress)

	option.AddFlag("no-progress", "Disable progress spinners")

	return option
}

func newNonInteractive(settings *Settings) *option.Bool {
	option := option.NewBoolOption(&settings.NonInteractive)

	option.AddFlag("non-interactive", "Disable interactive terminal")

	return option
}

func newToken(settings *Settings) *option.String {
	help := "Get yours from: https://environments.bunnyshell.com/access-token"

	option := option.NewStringOption(&settings.Profile.Token)

	option.AddFlagWithExtraHelp("token", "Authentication Token", help)

	return option
}

func newHost(settings *Settings) *option.String {
	option := option.NewStringOption(&settings.Profile.Host)

	flag := option.AddFlag("host", "Bunnyshell API Host")
	flag.Hidden = true

	return option
}

func newProfileName(settings *Settings) *option.String {
	help := fmt.Sprintf(
		`Local profile name. Find available profiles with "%s configure profile list"`,
		build.Name,
	)

	option := option.NewStringOption(&settings.Profile.Name)

	option.AddFlagWithExtraHelp("profile", "Use profile from config file", help)

	return option
}

func newOrganization(settings *Settings) *option.String {
	help := fmt.Sprintf(
		`Find available Organizations with "%s organization list"`,
		build.Name,
	)

	option := option.NewStringOption(&settings.Profile.Context.Organization)

	option.AddFlagWithExtraHelp("organization", "Filter by Organization", help)
	option.AddFlagWithExtraHelp("id", "Organization ID", help)

	return option
}

func newProject(settings *Settings) *option.String {
	help := fmt.Sprintf(
		`Find available Projects with "%s project list"`,
		build.Name,
	)

	option := option.NewStringOption(&settings.Profile.Context.Project)

	option.AddFlagWithExtraHelp("project", "Filter by Project", help)
	option.AddFlagWithExtraHelp("id", "Project ID", help)

	return option
}

func newEnvironment(settings *Settings) *option.String {
	help := fmt.Sprintf(
		`Find available Environments with "%s environment list"`,
		build.Name,
	)

	option := option.NewStringOption(&settings.Profile.Context.Environment)

	option.AddFlagWithExtraHelp("environment", "Filter by Environment", help)
	option.AddFlagWithExtraHelp("id", "Environment ID", help)

	return option
}

func newServiceComponent(settings *Settings) *option.String {
	help := fmt.Sprintf(
		`Find available Components with "%s components list"`,
		build.Name,
	)

	option := option.NewStringOption(&settings.Profile.Context.ServiceComponent)

	option.AddFlagWithExtraHelp("component", "Filter by ServiceComponent", help)
	option.AddFlagWithExtraHelp("id", "ServiceComponent ID", help)

	return option
}
