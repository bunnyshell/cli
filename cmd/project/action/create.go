package action

import (
	"errors"
	"fmt"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/api/project"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/util"
	"github.com/spf13/cobra"
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	createOptions := project.NewCreateOptions()

	command := &cobra.Command{
		Use: "create",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			createOptions.Organization = settings.Profile.Context.Organization

			fmt.Printf("%v", settings.Profile.Context)

			model, err := project.Create(createOptions)
			if err != nil {
				var apiError api.Error

				if errors.As(err, &apiError) {
					return handleCreateErrors(cmd, apiError, createOptions)
				}

				return lib.FormatCommandError(cmd, err)
			}

			return lib.FormatCommandData(cmd, model)
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.Organization.AddFlagWithExtraHelp(
		"organization",
		"Organization for the project",
		"Organizations contain projects along with build settings and project variables",
		util.FlagRequired,
	))

	createOptions.UpdateFlagSet(flags)

	mainCmd.AddCommand(command)
}

func handleCreateErrors(cmd *cobra.Command, apiError api.Error, createOptions *project.CreateOptions) error {
	if len(apiError.Violations) == 0 {
		return apiError
	}

	for _, violation := range apiError.Violations {
		cmd.Printf("Problem with creation: %s\n", violation.GetMessage())
	}

	return lib.ErrGeneric
}
