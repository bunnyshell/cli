package k8sIntegration

import (
	"bunnyshell.com/cli/pkg/api/k8s"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	listOptions := k8s.NewListOptions()

	command := &cobra.Command{
		Use:     "list",
		GroupID: mainGroup.ID,

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			listOptions.Environment = settings.Profile.Context.Environment
			listOptions.Organization = settings.Profile.Context.Organization

			return lib.ShowCollectionNoResponse(cmd, listOptions.Page, func(page int32) (lib.ModelWithPagination, error) {
				listOptions.Page = page

				return k8s.List(listOptions)
			})
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.Organization.GetFlag("organization"))
	flags.AddFlag(options.Environment.GetFlag("environment"))

	flags.Int32Var(&listOptions.Page, "page", listOptions.Page, "Listing Page")

	flags.StringVar(&listOptions.CloudProvider, "cloudProvider", listOptions.CloudProvider, "Filter by Cloud Provider")
	flags.StringVar(&listOptions.Status, "stauts", listOptions.Status, "Filter by Cloud Status")

	mainCmd.AddCommand(command)
}
