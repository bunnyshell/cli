package formatter

import (
	"fmt"
	"text/tabwriter"

	"bunnyshell.com/sdk"
)

func tabulateKubernetesCollection(w *tabwriter.Writer, data *sdk.PaginatedKubernetesIntegrationCollection) {
	fmt.Fprintf(w, "%v\t %v\t %v\t %v\t %v\t %v\n", "K8SIntegrationID", "OrganizationID", "Cloud", "Provider", "Cluster", "Status")

	if data.Embedded != nil {
		for _, item := range data.Embedded.Item {
			fmt.Fprintf(
				w,
				"%v\t %v\t %v\t %v\t %v\t %v\n",
				item.GetId(),
				item.GetOrganization(),
				item.GetCloudName(),
				item.GetCloudProvider(),
				item.GetClusterName(),
				item.GetStatus(),
			)
		}
	}
}

func tabulateKubernetesItem(w *tabwriter.Writer, item *sdk.KubernetesIntegrationItem) {
	fmt.Fprintf(w, "%v\t %v\n", "K8SIntegrationID", item.GetId())
	fmt.Fprintf(w, "%v\t %v\n", "OrganizationID", item.GetOrganization())
	fmt.Fprintf(w, "%v\t %v\n", "Cloud", item.GetCloudName())
	fmt.Fprintf(w, "%v\t %v\n", "Provider", item.GetCloudProvider())
	fmt.Fprintf(w, "%v\t %v\n", "Cluster", item.GetClusterName())
	fmt.Fprintf(w, "%v\t %v\n", "Status", item.GetStatus())
}
