package formatter

import (
	"fmt"
	"text/tabwriter"
	"time"

	"bunnyshell.com/sdk"
)

func tabulatePipelineCollection(writer *tabwriter.Writer, data *sdk.PaginatedPipelineCollection) {
	fmt.Fprintf(writer, "%v\t %v\t %v\t %v\t %v\n", "PipelineID", "EnvironmentID", "OrganizationID", "Description", "Status")

	if data.Embedded != nil {
		for _, item := range data.Embedded.Item {
			fmt.Fprintf(
				writer,
				"%v\t %v\t %v\t %v\t %v\n",
				item.GetId(),
				item.GetEnvironment(),
				item.GetOrganization(),
				item.GetDescription(),
				item.GetStatus(),
			)
		}
	}
}

func tabulatePipelineItem(writer *tabwriter.Writer, item *sdk.PipelineItem) {
	hasWebUrl := item.GetWebUrl() != ""

	fmt.Fprintf(writer, "%v\t %v\n", "PipelineID", item.GetId())
	fmt.Fprintf(writer, "%v\t %v\n", "EnvironmentID", item.GetEnvironment())
	fmt.Fprintf(writer, "%v\t %v\n", "OrganizationID", item.GetOrganization())
	fmt.Fprintf(writer, "%v\t %v\n", "Description", item.GetDescription())
	fmt.Fprintf(writer, "%v\t %v\n", "Status", item.GetStatus())

	if hasWebUrl {
		fmt.Fprintf(writer, "%v\t %v\n", "URL", item.GetWebUrl())
	}

	for index, stage := range item.GetStages() {
		if index == 0 {
			fmt.Fprintf(writer, "\n")
			fmt.Fprintf(writer, "%v\t %v\t %v\t %v\t %v\t %v\n", "Stages", "Name", "Duration", "Jobs", "JobsDone", "Status")
		}

		fmt.Fprintf(
			writer,
			"\t %v\t %v\t %v\t %v\t %v\n",
			stage.GetName(),
			time.Duration(stage.GetDuration())*time.Second,
			stage.GetJobsCount(),
			stage.GetCompletedJobsCount(),
			stage.GetStatus(),
		)
	}
}
