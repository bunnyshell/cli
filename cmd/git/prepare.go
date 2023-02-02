package git

import (
	"bunnyshell.com/cli/pkg/api/component/git"
	"bunnyshell.com/cli/pkg/config"
	gitHelper "bunnyshell.com/cli/pkg/helper/git"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	aggOptions := git.NewAggregateOptions()
	prepareOptions := gitHelper.NewPrepareOptions()

	command := &cobra.Command{
		Use: "prepare",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			aggOptions.Organization = settings.Profile.Context.Organization
			aggOptions.Project = settings.Profile.Context.Project
			aggOptions.Environment = settings.Profile.Context.Environment

			components, err := git.Aggregate(aggOptions)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			if err = gitHelper.PrintPrepareInfo(components, prepareOptions); err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			return nil
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.Organization.GetFlag("organization"))
	flags.AddFlag(options.Project.GetFlag("project"))
	flags.AddFlag(options.Environment.GetFlag("environment"))

	aggOptions.UpdateFlagSet(flags)
	prepareOptions.UpdateFlagSet(flags)

	mainCmd.AddCommand(command)
}
