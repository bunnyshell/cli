package profile

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/interactive"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/util"
	"bunnyshell.com/sdk"
	"github.com/spf13/cobra"
)

var tokenFormat = regexp.MustCompile(`^\d+:[0-9a-zA-z]{32}$`)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	profile := &settings.Profile
	timeout := settings.Timeout
	newProfileName := ""
	asDefaultProfile := false

	command := &cobra.Command{
		Use: "add",

		ValidArgsFunction: cobra.NoFileCompletions,

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := util.PersistentPreRunChain(cmd, args); err != nil {
				if errors.Is(err, config.ErrUnknownProfile) {
					return nil
				}
			}

			if errors.Is(config.MainManager.Error, config.ErrConfigLoad) {
				return nil
			}

			if !config.MainManager.HasProfile(newProfileName) {
				return nil
			}

			return config.ErrDuplicateProfile
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			if newProfileName == "" {
				if settings.NonInteractive {
					return interactive.ErrRequiredValue
				}

				profileName, err := interactive.Ask("Name:", getProfileNameValidator())
				if err != nil {
					return err
				}

				newProfileName = profileName
			} else {
				if err := getProfileNameValidator()(newProfileName); err != nil {
					return err
				}
			}

			settings.Timeout = 0 * time.Second

			for {
				if err := ensureToken(profile); err != nil {
					return err
				}

				organizations, resp, err := getOrganizations(settings.Profile)
				if err != nil {
					_ = lib.FormatCommandError(cmd, err)
					profile.Token = ""

					continue
				}

				if organizations.Embedded == nil || len(organizations.Embedded.Item) == 0 {
					return fmt.Errorf("create an organization in: %s", resp.Request.Host)
				}

				if err = setOrganization(&settings.Profile.Context, organizations.Embedded.Item); err != nil {
					if errors.Is(err, interactive.ErrNonInteractive) {
						return nil
					}

					return err
				}

				break
			}

			if err := config.MainManager.AddProfile(settings.Profile); err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			settings.Timeout = timeout

			return nil
		},

		PostRunE: func(cmd *cobra.Command, args []string) error {
			if settings.Verbosity != 0 {
				_ = lib.FormatCommandData(cmd, map[string]interface{}{
					"message": "Saved config file",
					"data":    settings.ConfigFile,
				})
			}

			if asDefaultProfile || askForDefault(cmd) {
				root := mainCmd.Root()
				root.SetArgs([]string{"configure", "profiles", "default", "--name", settings.Profile.Name})

				return root.Execute()
			}

			return nil
		},
	}

	config.MainManager.CommandWithAPI(command)

	flags := command.Flags()

	newProfileNameFlagName := "name"
	flags.StringVar(&newProfileName, newProfileNameFlagName, newProfileName, "Unique name for the new profile")
	_ = command.MarkFlagRequired(newProfileNameFlagName)

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
	)
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

func setOrganization(profileContext *config.Context, organizations []sdk.OrganizationCollection) error {
	if profileContext.Organization != "" {
		for _, organization := range organizations {
			if organization.Id == &profileContext.Organization {
				return nil
			}
		}

		return fmt.Errorf("%w: unknown organization (%s)", interactive.ErrInvalidValue, profileContext.Organization)
	}

	if config.GetSettings().NonInteractive {
		return interactive.ErrNonInteractive
	}

	index, _, err := interactive.Choose("Select Organization (empty to skip)", getOrganizationNames(organizations))
	profileContext.Organization = *organizations[index].Id

	return err
}

func getOrganizationNames(organizations []sdk.OrganizationCollection) []string {
	result := []string{}

	for _, organization := range organizations {
		result = append(result, *organization.Name)
	}

	return result
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

func getOrganizations(profile config.Profile) (*sdk.PaginatedOrganizationCollection, *http.Response, error) {
	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).OrganizationApi.OrganizationList(ctx)

	return request.Execute()
}
