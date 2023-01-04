package lib

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type CLI struct {
	Debug        bool
	NoProgress   bool
	ProfileName  string
	Profile      Profile
	OutputFormat string
	ConfigFile   string
	Verbosity    int
	Timeout      time.Duration

	commandFlags []*pflag.Flag
}

var defaultFormat = "stylish"
var defaultTimeout = 30 * time.Second

var CLIContext = CLI{}

func (c *CLI) Load(config Config) error {
	if !c.Debug {
		c.Debug = config.Debug
	}

	if c.OutputFormat == "" {
		c.OutputFormat = config.OutputFormat
	}

	switch c.OutputFormat {
	case "stylish", "json", "yaml", "yml":
	default:
		c.OutputFormat = defaultFormat
	}

	if c.Timeout == 0 {
		if config.Timeout != 0 {
			c.Timeout = config.Timeout
		} else {
			c.Timeout = defaultTimeout
		}
	}

	profileName := c.ProfileName
	if profileName == "" {
		profileName = config.DefaultProfile
	}

	profile, ok := config.Profiles[profileName]
	if !ok {
		return fmt.Errorf("profile %s not found", profileName)
	}

	if c.Profile.Token == "" {
		c.Profile.Token = profile.Token
		if c.Profile.Token != "" {
			c.markChangedToken()
		}
	}

	if c.Profile.Token == "" {
		token, ok := os.LookupEnv("BUNNYSHELL_TOKEN")
		if ok {
			c.Profile.Token = token
			c.markChangedToken()
		}
	}

	if c.Profile.Host == "" {
		c.Profile.Host = profile.Host
	}

	if c.Profile.Context.Organization == "" {
		c.Profile.Context.Organization = profile.Context.Organization
	}

	if c.Profile.Context.Project == "" {
		c.Profile.Context.Project = profile.Context.Project
	}

	if c.Profile.Context.Environment == "" {
		c.Profile.Context.Environment = profile.Context.Environment
	}

	if c.Profile.Context.ServiceComponent == "" {
		c.Profile.Context.ServiceComponent = profile.Context.ServiceComponent
	}

	return nil
}

// otherwise token would still be "required" and error out
func (c *CLI) markChangedToken() {
	for _, commandFlag := range c.commandFlags {
		commandFlag.Changed = true
	}
}

func (c *CLI) RequireTokenOnCommand(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&c.Profile.Token, "token", c.Profile.Token, "Authentication Token")
	c.commandFlags = append(c.commandFlags, cmd.PersistentFlags().Lookup("token"))

	cmd.MarkPersistentFlagRequired("token")
}

func (c *CLI) SetGlobalFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&CLIContext.ConfigFile, "configFile", "c", CLIContext.ConfigFile, "Config file")
	cmd.MarkPersistentFlagFilename("configFile", "yaml", "json")

	cmd.PersistentFlags().StringVarP(&CLIContext.ProfileName, "profile", "p", CLIContext.ProfileName, "Force profile usage from config file")
	cmd.PersistentFlags().StringVarP(&CLIContext.OutputFormat, "output", "o", CLIContext.OutputFormat, "Output format: stylish | json | yaml")
	cmd.PersistentFlags().BoolVarP(&CLIContext.Debug, "debug", "d", CLIContext.Debug, "Show network debug")
	cmd.PersistentFlags().BoolVar(&CLIContext.NoProgress, "no-progress", CLIContext.NoProgress, "Disable progress spinners")
	cmd.PersistentFlags().CountVarP(&CLIContext.Verbosity, "verbose", "v", "Number for the log level verbosity")

	cmd.PersistentFlags().DurationVarP(&CLIContext.Timeout, "timeout", "t", CLIContext.Timeout, "Network timeout on requests")
}

func MakeDefaultContext() {
	CLIContext.OutputFormat = defaultFormat
	CLIContext.Timeout = defaultTimeout
}

func LoadViperConfigIntoContext() {
	if err := viper.ReadInConfig(); err != nil {
		if CLIContext.ConfigFile != "" {
			fmt.Fprintln(os.Stderr, "[LoadConfigError]", err)
		}
		if CLIContext.Verbosity != 0 {
			fmt.Fprintln(os.Stderr, "[LoadConfigError]", err)
		}
		return
	}

	config, err := GetConfig()
	if err != nil {
		if CLIContext.Verbosity != 0 {
			fmt.Fprintln(os.Stderr, "[LoadConfigError]", err)
		}
		return
	}

	err = CLIContext.Load(*config)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[LoadConfigError]", err)
	}
}
