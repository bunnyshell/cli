package profile

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/util"
	"bunnyshell.com/sdk"
)

func init() {
	var profile lib.Profile
	var profileName string

	var addProfileCommand = &cobra.Command{
		Use: "add name [token organization project]",

		ValidArgsFunction: cobra.NoFileCompletions,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			var err error

			if profileName == "" {
				profileName, err = util.Ask("Name:", getProfileNameValidator())
			} else {
				err = getProfileNameValidator()(profileName)
			}

			if err != nil {
				return err
			}

			lib.CLIContext.Timeout = 0 * time.Second
			for {
				if err := ensureToken(&profile); err != nil {
					return err
				}

				organizations, r, err := getOrganizations(profile)
				if err != nil {
					lib.FormatCommandError(cmd, err)
					profile.Token = ""
					continue
				}

				if organizations.Embedded == nil || len(organizations.Embedded.Item) == 0 {
					return fmt.Errorf("create an organization in: %s", r.Request.Host)
				}

				ok, err := util.Confirm("Set a default organization")
				if err != nil {
					return err
				}

				if ok {
					if err := setOrganization(&profile, organizations.Embedded.Item); err != nil {
						return err
					}
				}

				break
			}

			err = lib.AddProfile(profile, profileName)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			if err := viper.WriteConfig(); err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if lib.CLIContext.Verbosity != 0 {
				lib.FormatCommandData(cmd, map[string]interface{}{
					"message": "Saved config file",
					"data":    viper.ConfigFileUsed(),
				})
			}

			ok, err := util.Confirm("Set as default profile?")
			if err != nil {
				cmd.PrintErr("Could not determine user input", err)
				return
			}

			if !ok {
				return
			}

			root := mainCmd.Root()
			root.SetArgs([]string{"configure", "profiles", "default", "--name", profileName})
			root.Execute()
		},
	}

	addProfileCommand.Flags().StringVar(&profile.Host, "host", profile.Host, "Host for the new profile")
	addProfileCommand.Flags().StringVar(&profile.Token, "token", profile.Token, "Token for the new profile")
	addProfileCommand.Flags().StringVar(&profile.Context.Organization, "organization", profile.Context.Organization, "AutoFilter on Organization")
	addProfileCommand.Flags().StringVar(&profile.Context.Project, "project", profile.Context.Project, "AutoFilter on Project")
	addProfileCommand.Flags().StringVar(&profile.Context.Environment, "environment", profile.Context.Environment, "AutoFilter on Enviroment")
	addProfileCommand.Flags().StringVar(&profile.Context.ServiceComponent, "serviceComponent", profile.Context.ServiceComponent, "AutoFilter on ServiceComponent")
	addProfileCommand.Flags().StringVar(&profileName, "name", profileName, "Unique name for the new profile")

	mainCmd.AddCommand(addProfileCommand)
}

func getProfileNameValidator() func(interface{}) error {
	return util.All(
		util.Lowercase(),
		util.AssertMinimumLength(4),
	)
}

func ensureToken(profile *lib.Profile) error {
	if profile.Token != "" {
		return nil
	}

	token, err := util.AskSecretWithHelp("Token:", "Get yours from: https://environments.bunnyshell.com/access-token", validateToken)
	if err != nil {
		return err
	}

	profile.Token = token

	return nil
}

func setOrganization(profile *lib.Profile, organizations []sdk.OrganizationCollection) error {
	if profile.Context.Organization != "" {
		return nil
	}

	index, _, err := util.Choose("Select Organization", getOrganizationNames(organizations))
	profile.Context.Organization = *organizations[index].Id

	return err
}
func getOrganizationNames(organizations []sdk.OrganizationCollection) []string {
	var result []string
	for _, organization := range organizations {
		result = append(result, *organization.Name)
	}
	return result
}

func validateToken(input interface{}) error {
	chunks := strings.Split(input.(string), ":")
	if len(chunks) != 2 {
		return errors.New("invalid token detected")
	}

	if len(chunks[1]) != 32 {
		return errors.New("invalid token detected")
	}

	if _, err := strconv.Atoi(chunks[0]); err != nil {
		return errors.New("invalid token detected")
	}

	return nil
}

func getOrganizations(profile lib.Profile) (*sdk.PaginatedOrganizationCollection, *http.Response, error) {
	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetApiFromProfile(profile).OrganizationApi.OrganizationList(ctx)

	return request.Execute()
}
