package formatter

import (
	"fmt"
	"text/tabwriter"
	"time"

	"bunnyshell.com/sdk"
)

func formatStartedAt(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format(time.RFC3339)
}

func tabulateWorkflowCollection(writer *tabwriter.Writer, data *sdk.PaginatedWorkflowCollection) {
	fmt.Fprintf(writer, "%v\t %v\t %v\t %v\t %v\t %v\t %v\n", "PipelineID", "EventID", "EnvironmentID", "OrganizationID", "Description", "Status", "StartedAt")

	if data.Embedded != nil {
		for _, item := range data.Embedded.Item {
			fmt.Fprintf(
				writer,
				"%v\t %v\t %v\t %v\t %v\t %v\t %v\n",
				item.GetId(),
				item.GetEvent(),
				item.GetEnvironment(),
				item.GetOrganization(),
				item.GetDescription(),
				item.GetStatus(),
				formatStartedAt(item.GetStartedAt()),
			)
		}
	}
}

func tabulateWorkflowItem(writer *tabwriter.Writer, item *sdk.WorkflowItem) {
	hasWebUrl := item.GetWebUrl() != ""

	fmt.Fprintf(writer, "%v\t %v\n", "PipelineID", item.GetId())
	fmt.Fprintf(writer, "%v\t %v\n", "EnvironmentID", item.GetEnvironment())
	fmt.Fprintf(writer, "%v\t %v\n", "OrganizationID", item.GetOrganization())
	fmt.Fprintf(writer, "%v\t %v\n", "EventID", item.GetEvent())
	fmt.Fprintf(writer, "%v\t %v\n", "Description", item.GetDescription())
	fmt.Fprintf(writer, "%v\t %v\n", "Status", item.GetStatus())
	fmt.Fprintf(writer, "%v\t %v\n", "StartedAt", formatStartedAt(item.GetStartedAt()))
	fmt.Fprintf(writer, "%v\t %v\n", "Jobs", fmt.Sprintf("%d/%d completed", item.GetCompletedJobsCount(), item.GetJobsCount()))
	fmt.Fprintf(writer, "%v\t %v\n", "Duration", time.Duration(item.GetDuration())*time.Second)

	if hasWebUrl {
		fmt.Fprintf(writer, "%v\t %v\n", "URL", item.GetWebUrl())
	}
}
