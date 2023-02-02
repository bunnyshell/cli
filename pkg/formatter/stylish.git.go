package formatter

import (
	"fmt"
	"text/tabwriter"

	"bunnyshell.com/sdk"
)

func tabulateComponentGitCollection(writer *tabwriter.Writer, data *sdk.PaginatedComponentGitCollection) {
	fmt.Fprintf(writer, "%v\t %v\t %v\t %v\t %v\t %v\t %v\n", "EnvironmentID", "ComponentID", "Name", "Repository", "Branch", "Sha", "DeployedSha")

	if !data.HasEmbedded() {
		return
	}

	for _, item := range data.Embedded.Item {
		fmt.Fprintf(
			writer,
			"%v\t %v\t %v\t %v\t %v\t %v\t %v\n",
			item.GetEnvironment(),
			item.GetId(),
			item.GetName(),
			item.GetRepository(),
			item.GetRefName(),
			item.GetRefSha(),
			item.GetDeployedSha(),
		)
	}
}

func tabulateComponentGitItem(writer *tabwriter.Writer, item *sdk.ComponentGitItem) {
	fmt.Fprintf(writer, "%v\t %v\n", "EnvironmentID", item.GetEnvironment())
	fmt.Fprintf(writer, "%v\t %v\n", "ComponentID", item.GetId())
	fmt.Fprintf(writer, "%v\t %v\n", "Name", item.GetName())
	fmt.Fprintf(writer, "%v\t %v\n", "Repository", item.GetName())
	fmt.Fprintf(writer, "%v\t %v\n", "Branch", item.GetRefName())
	fmt.Fprintf(writer, "%v\t %v\n", "Sha", item.GetRefSha())
	fmt.Fprintf(writer, "%v\t %v\n", "DeployedSha", item.GetDeployedSha())
}
