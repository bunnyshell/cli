package environment

import (
	"bunnyshell.com/cli/pkg/api/environment"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/util"
	"github.com/spf13/cobra"
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	listOptions := environment.NewListOptions()

	command := &cobra.Command{
		Use:     "list",
		GroupID: mainGroup.ID,

		ValidArgsFunction: cobra.NoFileCompletions,

		PersistentPreRunE: util.PersistentPreRunChain,

		RunE: func(cmd *cobra.Command, args []string) error {
			listOptions.Organization = settings.Profile.Context.Organization
			listOptions.Project = settings.Profile.Context.Project

			return lib.ShowCollectionNoResponse(cmd, listOptions.Page, func(page int32) (lib.ModelWithPagination, error) {
				listOptions.Page = page

				return environment.List(listOptions)
			})
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.Organization.GetFlag("organization"))
	flags.AddFlag(options.Project.GetFlag("project"))

	flags.StringVar(&listOptions.KubernetesIntegration, "k8s-cluster", listOptions.KubernetesIntegration, "Filter by K8s Cluster")

	flags.Int32Var(&listOptions.Page, "page", listOptions.Page, "Listing Page")

	flags.StringVar(&listOptions.Type, "type", listOptions.Type, "Filter by Type")
	flags.StringVar(&listOptions.ClusterStatus, "clusterStatus", listOptions.ClusterStatus, "Filter by Cluster Status")
	flags.StringVar(&listOptions.OperationStatus, "operationStatus", listOptions.OperationStatus, "Filter by Operation Status")

	mainCmd.AddCommand(command)
}
