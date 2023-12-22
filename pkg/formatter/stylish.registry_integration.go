package formatter

import (
	"fmt"
	"text/tabwriter"

	"bunnyshell.com/sdk"
)

func tabulateRegistryIntegrationsCollection(w *tabwriter.Writer, data *sdk.PaginatedRegistryIntegrationCollection) {
	fmt.Fprintf(w, "%v\t %v\t %v\t %v\t %v\n", "IntegrationID", "OrganizationID", "Name", "Provider", "Status")

	if data.Embedded != nil {
		for _, item := range data.Embedded.Item {
			fmt.Fprintf(
				w,
				"%v\t %v\t %v\t %v\t %v\n",
				item.GetId(),
				item.GetOrganization(),
				item.GetName(),
				item.GetProviderName(),
				item.GetStatus(),
			)
		}
	}
}

func tabulateRegistryIntegrationItem(w *tabwriter.Writer, item *sdk.RegistryIntegrationItem) {
	fmt.Fprintf(w, "%v\t %v\n", "IntegrationID", item.GetId())
	fmt.Fprintf(w, "%v\t %v\n", "OrganizationID", item.GetOrganization())
	fmt.Fprintf(w, "%v\t %v\n", "Name", item.GetName())
	fmt.Fprintf(w, "%v\t %v\n", "Provider", item.GetProviderName())
	fmt.Fprintf(w, "%v\t %v\n", "Status", item.GetStatus())
}
