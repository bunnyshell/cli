package profile

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"bunnyshell.com/cli/pkg/lib"
)

func init() {
	profileName := &lib.CLIContext.ProfileName
	var organization string
	var project string
	var environment string
	var serviceComponent string

	var contextCommand = &cobra.Command{
		Use: "context",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			profile, err := lib.GetProfile(lib.CLIContext.ProfileName)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			var updates []string
			if organization != "" {
				updates = append(updates, "organization")
				profile.Context.Organization = organization
			}

			if project != "" {
				updates = append(updates, "project")
				profile.Context.Project = project
			}

			if environment != "" {
				updates = append(updates, "environment")
				profile.Context.Environment = environment
			}

			if serviceComponent != "" {
				updates = append(updates, "service component")
				profile.Context.ServiceComponent = serviceComponent
			}

			if len(updates) == 0 {
				return errors.New("no context changes")
			}

			if err := viper.WriteConfig(); err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			return lib.FormatCommandData(cmd, map[string]interface{}{
				"message": "Updated profile context",
				"data":    updates,
			})
		},
	}

	contextCommand.Flags().StringVar(profileName, "name", *profileName, "Name of the profile")
	contextCommand.MarkFlagRequired("name")

	contextCommand.Flags().StringVar(&organization, "organization", organization, "Set Organization context for all resources")
	contextCommand.Flags().StringVar(&project, "project", project, "Set Project context for all resources")
	contextCommand.Flags().StringVar(&environment, "environment", environment, "Set Organization context for all resources")
	contextCommand.Flags().StringVar(&serviceComponent, "serviceComponent", serviceComponent, "Set Organization context for all resources")

	mainCmd.AddCommand(contextCommand)
}
