package profile

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"bunnyshell.com/cli/pkg/api/organization"
	"bunnyshell.com/cli/pkg/build"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/config/option"
	"bunnyshell.com/cli/pkg/interactive"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/util"
	"github.com/spf13/cobra"
)

var (
	tokenFormat = regexp.MustCompile(`^\d+:[0-9a-zA-z]{32}$`)

	ErrInvalidToken = errors.New("invalid token")
)

func init() {
	settings := config.GetSettings()

	profile := &config.Profile{}

	asDefaultProfile := false

	command := &cobra.Command{
		Use: "add",

		ValidArgsFunction: cobra.NoFileCompletions,

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := util.PersistentPreRunChain(cmd, args); err != nil {
				if errors.Is(err, config.ErrUnknownProfile) {
					return nil
				}

				return err
			}

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ensureProfileName(profile); err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			settings.Timeout = 0 * time.Second

			for {
				if err := ensureToken(profile); err != nil {
					if errors.Is(err, interactive.ErrInvalidValue) {
						continue
					}

					return lib.FormatCommandError(cmd, err)
				}

				if err := checkToken(profile); err != nil {
					return lib.FormatCommandError(cmd, err)
				}

				break
			}

			if err := askToFillContextOrSkip(profile); err != nil {
				if errors.Is(err, interactive.ErrNonInteractive) {
					return nil
				}

				return lib.FormatCommandError(cmd, err)
			}

			if err := config.MainManager.AddProfile(*profile); err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			if asDefaultProfile || askForDefault(cmd) {
				if err := setDefaultProfile(profile); err != nil {
					return lib.FormatCommandData(cmd, err)
				}
			}

			if err := config.MainManager.Save(); err != nil {
				return lib.FormatCommandData(cmd, err)
			}

			return nil
		},
	}

	flags := command.Flags()

	flags.AddFlag(getNewProfileNameOption(&profile.Name).GetRequiredFlag("name"))

	flags.StringVar(
		&profile.Token,
		"token",
		profile.Token,
		"Set API token for the new profile",
	)

	flags.StringVar(
		&profile.Host,
		"host",
		profile.Host,
		"Set API host for the new profile",
	)

	_ = flags.MarkHidden("host")

	flags.StringVar(
		&profile.Context.Organization,
		"organization",
		profile.Context.Organization,
		"Set Organization context for all resources",
	)
	flags.StringVar(
		&profile.Context.Project,
		"project",
		profile.Context.Project,
		"Set Project context for all resources",
	)
	flags.StringVar(
		&profile.Context.Environment,
		"environment",
		profile.Context.Environment,
		"Set Environment context for all resources",
	)
	flags.StringVar(
		&profile.Context.ServiceComponent,
		"service",
		profile.Context.ServiceComponent,
		"Set ServiceComponent context for all resources",
	)

	flags.BoolVar(&asDefaultProfile, "default", asDefaultProfile, "Set as default profile")

	mainCmd.AddCommand(command)
}

func getNewProfileNameOption(value *string) *option.String {
	usage := "Unique name for the new profile"
	help := usage

	idOption := option.NewStringOption(value)

	idOption.AddFlagWithExtraHelp("name", usage, help)

	return idOption
}

func askForDefault(command *cobra.Command) bool {
	if config.GetSettings().NonInteractive {
		return false
	}

	setAsDefault, err := interactive.ConfirmWithHelp(
		"Set as default profile?",
		"The default profile is automatically used when running commands without the --profile flag\n"+
			fmt.Sprintf("You can change the default profile with '%s configure profiles default' command\n", build.Name)+
			"See more at https://documentation.bunnyshell.com/docs/bunnyshell-cli-authentication#create-a-profile",
	)
	if err != nil {
		command.PrintErr("Could not determine user input", err)

		return false
	}

	return setAsDefault
}

func getProfileNameValidator() func(interface{}) error {
	return interactive.All(
		interactive.Lowercase(),
		interactive.AssertMinimumLength(4),
		func(input interface{}) error {
			str, ok := input.(string)
			if !ok {
				return interactive.ErrInvalidValue
			}

			if config.MainManager.HasProfile(str) {
				return config.ErrDuplicateProfile
			}

			return nil
		},
	)
}

func ensureProfileName(profile *config.Profile) error {
	if profile.Name != "" {
		return getProfileNameValidator()(profile.Name)
	}

	if config.GetSettings().NonInteractive {
		return interactive.ErrRequiredValue
	}

	profileName, err := interactive.Ask("Name:", getProfileNameValidator())
	if err != nil {
		return err
	}

	profile.Name = profileName

	return nil
}

func ensureToken(profile *config.Profile) error {
	if profile.Token != "" {
		if err := validateToken(profile.Token); err != nil {
			profile.Token = ""

			return ErrInvalidToken
		}

		return nil
	}

	if config.GetSettings().NonInteractive {
		return fmt.Errorf("%w (token)", interactive.ErrRequiredValue)
	}

	tokenFlag := config.GetOptions().Token.GetMainFlag()
	help := util.GetHelp(tokenFlag)

	token, err := interactive.AskSecretWithHelp("Token:", help, validateToken)
	if err != nil {
		return err
	}

	profile.Token = token

	return nil
}

func checkToken(profile *config.Profile) error {
	listOptions := organization.NewListOptions()
	listOptions.Profile = profile

	_, err := organization.List(listOptions)

	return err
}

func validateToken(input interface{}) error {
	value, ok := input.(string)
	if !ok {
		return interactive.ErrInvalidValue
	}

	if !tokenFormat.Match([]byte(value)) {
		return fmt.Errorf("%w: token is invalid", interactive.ErrInvalidValue)
	}

	return nil
}
