package formatter

import (
	"fmt"
	"text/tabwriter"
	"time"

	"bunnyshell.com/sdk"
)

func tabulateWorkflowJobCollection(writer *tabwriter.Writer, data *sdk.PaginatedWorkflowJobCollection) {
	fmt.Fprintf(
		writer,
		"%v\t %v\t %v\t %v\t %v\t %v\t %v\t %v\n",
		"JobID",
		"PipelineID",
		"Name",
		"Type",
		"Status",
		"StartedAt",
		"Duration",
		"AllowedToFail",
	)

	if data.Embedded != nil {
		for _, item := range data.Embedded.Item {
			duration := ""
			if value, ok := item.GetDurationOk(); ok && value != nil {
				duration = (time.Duration(*value) * time.Second).String()
			}

			allowedToFail := ""
			if value, ok := item.GetAllowedToFailOk(); ok && value != nil {
				allowedToFail = fmt.Sprintf("%t", *value)
			}

			startedAt := ""
			if value, ok := item.GetStartedAtOk(); ok && value != nil {
				startedAt = value.Format(time.RFC3339)
			}

			fmt.Fprintf(
				writer,
				"%v\t %v\t %v\t %v\t %v\t %v\t %v\t %v\n",
				item.GetId(),
				item.GetWorkflow(),
				item.GetName(),
				item.GetType(),
				item.GetStatus(),
				startedAt,
				duration,
				allowedToFail,
			)
		}
	}
}
