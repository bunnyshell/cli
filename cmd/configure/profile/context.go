package profile

import (
	"errors"

	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

var errNoContextChanges = errors.New("no context changes")

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()
	profile := &settings.Profile
	updateContext := &profile.Context

	command := &cobra.Command{
		Use: "context",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			configProfile, _ := config.MainManager.GetProfile(settings.Profile.Name)
			cfgContext := &configProfile.Context

			var updates []string
			if cfgContext.Organization != updateContext.Organization {
				updates = append(updates, "organization")
				cfgContext.Organization = updateContext.Organization
			}

			if cfgContext.Project != updateContext.Project {
				updates = append(updates, "project")
				cfgContext.Project = updateContext.Project
			}

			if cfgContext.Environment != updateContext.Environment {
				updates = append(updates, "environment")
				cfgContext.Environment = updateContext.Environment
			}

			if cfgContext.ServiceComponent != updateContext.ServiceComponent {
				updates = append(updates, "service component")
				cfgContext.ServiceComponent = updateContext.ServiceComponent
			}

			if len(updates) == 0 {
				return errNoContextChanges
			}

			config.MainManager.SetProfile(*profile)

			if err := config.MainManager.Save(); err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			return lib.FormatCommandData(cmd, map[string]interface{}{
				"message": "Updated profile context",
				"data":    updates,
			})
		},
	}

	flags := command.Flags()

	profileNameFlag := options.ProfileName.CloneMainFlag()
	flags.AddFlag(profileNameFlag)
	_ = command.MarkFlagRequired(profileNameFlag.Name)

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

	mainCmd.AddCommand(command)
}
