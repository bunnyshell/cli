package profile

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"bunnyshell.com/cli/pkg/api/organization"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/interactive"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/util"
	"github.com/spf13/cobra"
)

var tokenFormat = regexp.MustCompile(`^\d+:[0-9a-zA-z]{32}$`)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	profile := &settings.Profile

	profileName := ""
	asDefaultProfile := false

	command := &cobra.Command{
		Use: "add",

		ValidArgsFunction: cobra.NoFileCompletions,

		PersistentPreRunE: util.PersistentPreRunChain,

		RunE: func(cmd *cobra.Command, args []string) error {
			profile.Name = profileName

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

			if err := askToFillContext(profile); err != nil {
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

	config.MainManager.CommandWithAPI(command)

	flags := command.Flags()

	profileNameFlagName := "name"
	flags.StringVar(&profileName, profileNameFlagName, profileName, "Unique name for the new profile")
	_ = command.MarkFlagRequired(profileNameFlagName)

	flags.AddFlag(
		options.Organization.AddFlag("organization", "Set Organization context for all resources"),
	)
	flags.AddFlag(
		options.Project.AddFlag("project", "Set Project context for all resources"),
	)
	flags.AddFlag(
		options.Environment.AddFlag("environment", "Set Environment context for all resources"),
	)
	flags.AddFlag(
		options.ServiceComponent.AddFlag("serviceComponent", "Set ServiceComponent context for all resources"),
	)

	flags.BoolVar(&asDefaultProfile, "default", asDefaultProfile, "Set as default profile")

	mainCmd.AddCommand(command)
}

func askForDefault(command *cobra.Command) bool {
	if config.GetSettings().NonInteractive {
		return false
	}

	setAsDefault, err := interactive.Confirm("Set as default profile?")
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

			return err
		}

		return nil
	}

	if config.GetSettings().NonInteractive {
		return fmt.Errorf("%w (token)", interactive.ErrRequiredValue)
	}

	help := "Get yours from: https://environments.bunnyshell.com/access-token"

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
