package formatter

import (
	"fmt"
	"text/tabwriter"

	"bunnyshell.com/sdk"
)

func tabulateAggregateEndpoint(writer *tabwriter.Writer, components []sdk.ComponentEndpointCollection) {
	if len(components) == 0 {
		fmt.Fprintln(writer, "Environment has no defined public endpoints")

		return
	}

	for index, item := range components {
		if index != 0 {
			fmt.Fprintln(writer)
		}

		fmt.Fprintf(writer, "%v\t %v\n", "EnvironmentID", item.GetEnvironment())
		fmt.Fprintf(writer, "%v\t %v\n", "ComponentID", item.GetId())
		fmt.Fprintf(writer, "%v\t %v\n", "Name", item.GetName())

		for index, endpoint := range item.GetEndpoints() {
			if index == 0 {
				fmt.Fprintf(writer, "%v\t %v\n", "Endpoints", endpoint)
			} else {
				fmt.Fprintf(writer, "\t %v\n", endpoint)
			}
		}
	}
}
