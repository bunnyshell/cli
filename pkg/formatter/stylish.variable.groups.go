package formatter

import (
	"fmt"
	"text/tabwriter"

	"bunnyshell.com/sdk"
)

func tabulateEnvironCollection(w *tabwriter.Writer, data *sdk.PaginatedEnvironItemCollection) {
	fmt.Fprintf(w, "%v\t %v\t %v\t %v\t %v\n", "EnvVarID", "EnvironmentID", "OrganizationID", "Group", "Name")

	if data.Embedded != nil {
		for _, item := range data.Embedded.Item {
			fmt.Fprintf(
				w,
				"%v\t %v\t %v\t %v\t %v\n",
				item.GetId(),
				item.GetEnvironment(),
				item.GetOrganization(),
				item.GetGroupName(),
				item.GetName(),
			)
		}
	}
}

func tabulateEnvironItem(w *tabwriter.Writer, item *sdk.EnvironItemItem) {
	fmt.Fprintf(w, "%v\t %v\n", "EnvironmentVariableID", item.GetId())
	fmt.Fprintf(w, "%v\t %v\n", "EnvironmentID", item.GetEnvironment())
	fmt.Fprintf(w, "%v\t %v\n", "OrganizationID", item.GetOrganization())
	fmt.Fprintf(w, "%v\t %v\n", "Group", item.GetGroupName())
	fmt.Fprintf(w, "%v\t %v\n", "Name", item.GetName())
	fmt.Fprintf(w, "%v\t %v\n", "Value", item.GetValue())
	fmt.Fprintf(w, "%v\t %v\n", "Secret", item.GetSecret())
}
