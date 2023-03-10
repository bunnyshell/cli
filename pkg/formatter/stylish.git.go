package formatter

import (
	"fmt"
	"text/tabwriter"

	"bunnyshell.com/sdk"
)

func tabulateComponentGitCollection(writer *tabwriter.Writer, data *sdk.PaginatedComponentGitCollection) {
	fmt.Fprintf(writer, "%v\t %v\t %v\t %v\t %v\t %v\t %v\t %v\n", "EnvironmentID", "ComponentID", "Name", "Repository", "Branch", "Path", "Sha", "DeployedSha")

	if !data.HasEmbedded() {
		return
	}

	for _, item := range data.Embedded.Item {
		fmt.Fprintf(
			writer,
			"%v\t %v\t %v\t %v\t %v\t %v\t %v\t %v\n",
			item.GetEnvironment(),
			item.GetId(),
			item.GetName(),
			item.GetRepository(),
			item.GetRefName(),
			item.GetPath(),
			item.GetRefSha(),
			item.GetDeployedSha(),
		)
	}
}

func tabulateComponentGitList(writer *tabwriter.Writer, data []sdk.ComponentGitCollection) {
	fmt.Fprintf(writer, "%v\t %v\t %v\t %v\t %v\t %v\n", "ComponentID", "Name", "Repository", "Branch", "Path", "Deployed")

	for _, item := range data {
		deployed := "latest"

		if item.GetRefSha() != item.GetDeployedSha() {
			deployed = "outdated"
		}

		fmt.Fprintf(
			writer,
			"%v\t %v\t %v\t %v\t %v\t %v\n",
			item.GetId(),
			item.GetName(),
			item.GetRepository(),
			item.GetRefName(),
			item.GetPath(),
			deployed,
		)
	}
}

func tabulateComponentGitItem(writer *tabwriter.Writer, item *sdk.ComponentGitItem) {
	fmt.Fprintf(writer, "%v\t %v\n", "EnvironmentID", item.GetEnvironment())
	fmt.Fprintf(writer, "%v\t %v\n", "ComponentID", item.GetId())
	fmt.Fprintf(writer, "%v\t %v\n", "Name", item.GetName())
	fmt.Fprintf(writer, "%v\t %v\n", "Repository", item.GetRepository())
	fmt.Fprintf(writer, "%v\t %v\n", "Branch", item.GetRefName())
	fmt.Fprintf(writer, "%v\t %v\n", "Path", item.GetPath())
	fmt.Fprintf(writer, "%v\t %v\n", "Sha", item.GetRefSha())
	fmt.Fprintf(writer, "%v\t %v\n", "DeployedSha", item.GetDeployedSha())
}
