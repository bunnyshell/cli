package formatter

import (
	"fmt"
	"text/tabwriter"

	"bunnyshell.com/sdk"
)

func tabulateTemplatesRepositoryCollection(writer *tabwriter.Writer, data *sdk.PaginatedTemplatesRepositoryCollection) {
	fmt.Fprintf(writer, "%v\t %v\t %v\t %v\t %v\t %v\n", "TemplateRepositoryID", "OrganizationID", "Name", "Repository", "Branch", "LastSyncSHA")

	if data.Embedded != nil {
		for _, item := range data.Embedded.Item {
			fmt.Fprintf(
				writer,
				"%v\t %v\t %v\t %v\t %v\t %v\n",
				item.GetId(),
				item.GetOrganization(),
				item.GetName(),
				item.GetGitRepositoryUrl(),
				item.GetGitRef(),
				item.GetLastSyncSha(),
			)
		}
	}
}

func tabulateTemplatesRepositoryItem(writer *tabwriter.Writer, item *sdk.TemplatesRepositoryItem) {
	fmt.Fprintf(writer, "%v\t %v\n", "TemplateRepositoryID", item.GetId())
	fmt.Fprintf(writer, "%v\t %v\n", "OrganizationID", item.GetOrganization())
	fmt.Fprintf(writer, "%v\t %v\n", "Name", item.GetName())
	fmt.Fprintf(writer, "%v\t %v\n", "Repository", item.GetGitRepositoryUrl())
	fmt.Fprintf(writer, "%v\t %v\n", "Branch", item.GetGitRef())
	fmt.Fprintf(writer, "%v\t %v\n", "LastSyncSHA", item.GetLastSyncSha())
}
