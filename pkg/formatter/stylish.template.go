package formatter

import (
	"fmt"
	"text/tabwriter"

	"bunnyshell.com/sdk"
)

func tabulateTemplateCollection(writer *tabwriter.Writer, data *sdk.PaginatedTemplateCollection) {
	fmt.Fprintf(writer, "%v\t %v\t %v\t %v\t %v\t %v\n", "TemplateID", "OrganizationID", "TemplatesRepositoryID", "Key", "Name", "Git SHA")

	if data.Embedded != nil {
		for _, item := range data.Embedded.Item {
			fmt.Fprintf(
				writer,
				"%v\t %v\t %v\t %v\t %v\t %v\n",
				item.GetId(),
				item.GetOrganization(),
				item.GetTemplatesRepository(),
				item.GetKey(),
				item.GetName(),
				item.GetGitSha(),
			)
		}
	}
}

func tabulateTemplateItem(writer *tabwriter.Writer, item *sdk.TemplateItem) {
	fmt.Fprintf(writer, "%v\t %v\n", "TemplateID", item.GetId())
	fmt.Fprintf(writer, "%v\t %v\n", "OrganizationID", item.GetOrganization())
	fmt.Fprintf(writer, "%v\t %v\n", "TemplatesRepositoryID", item.GetTemplatesRepository())
	fmt.Fprintf(writer, "%v\t %v\n", "Key", item.GetKey())
	fmt.Fprintf(writer, "%v\t %v\n", "Name", item.GetName())
	fmt.Fprintf(writer, "%v\t %v\n", "Short Description", item.GetShortDescription())
	fmt.Fprintf(writer, "%v\t %v\n", "Git SHA", item.GetGitSha())

	for index, tag := range item.GetTags() {
		if index == 0 {
			fmt.Fprintf(writer, "\n")
			fmt.Fprintf(writer, "%v\t %v\n", "Tags", tag)
		} else {
			fmt.Fprintf(writer, "\t %v\n", tag)
		}
	}

	tabulateTemplateVariableFromItem(writer, item)
}
