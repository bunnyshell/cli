package pipeline

import (
	"bunnyshell.com/cli/pkg/api/pipeline"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	listOptions := pipeline.NewListOptions()

	command := &cobra.Command{
		Use: "list",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			listOptions.Environment = settings.Profile.Context.Environment
			listOptions.Organization = settings.Profile.Context.Organization

			return lib.ShowCollectionNoResponse(cmd, listOptions.Page, func(page int32) (lib.ModelWithPagination, error) {
				listOptions.Page = page

				return pipeline.List(listOptions)
			})
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.Organization.GetFlag("organization"))
	flags.AddFlag(options.Environment.GetFlag("environment"))

	flags.StringVar(&listOptions.Event, "event", listOptions.Event, "Filter by event")
	flags.Int32Var(&listOptions.Page, "page", listOptions.Page, "Listing Page")

	mainCmd.AddCommand(command)
}
