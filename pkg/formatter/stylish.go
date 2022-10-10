package formatter

import (
	"bytes"
	"fmt"
	"text/tabwriter"

	"bunnyshell.com/sdk"
)

func stylish(data interface{}) ([]byte, error) {
	var b bytes.Buffer
	var err error = nil

	w := tabwriter.NewWriter(&b, 1, 1, 1, ' ', tabwriter.Debug)

	switch t := data.(type) {
	case *sdk.PaginatedOrganizationCollection:
		tabulateOrganizationCollection(w, t)
		tabulatePagination(w, t.GetPage(), t.GetItemsPerPage(), t.GetTotalItems())
	case *sdk.PaginatedProjectCollection:
		tabulateProjectCollection(w, t)
		tabulatePagination(w, t.GetPage(), t.GetItemsPerPage(), t.GetTotalItems())
	case *sdk.PaginatedEnvironmentCollection:
		tabulateEnvironmentCollection(w, t)
		tabulatePagination(w, t.GetPage(), t.GetItemsPerPage(), t.GetTotalItems())
	case *sdk.PaginatedComponentCollection:
		tabulateComponentCollection(w, t)
		tabulatePagination(w, t.GetPage(), t.GetItemsPerPage(), t.GetTotalItems())
	case *sdk.PaginatedEventCollection:
		tabulateEventCollection(w, t)
		tabulatePagination(w, t.GetPage(), t.GetItemsPerPage(), t.GetTotalItems())
	case *sdk.OrganizationItem:
		tabulateOrganizationItem(w, t)
	case *sdk.ProjectItem:
		tabulateProjectItem(w, t)
	case *sdk.EnvironmentItem:
		tabulateEnvironmentItem(w, t)
	case *sdk.ComponentItem:
		tabulateComponentItem(w, t)
	case *sdk.EventItem:
		tabulateEventItem(w, t)
	case *sdk.ProblemGeneric:
		tabulateGeneric(w, t)
	default:
		fmt.Fprintf(w, "JSON: ")
		var jsonBytes []byte
		jsonBytes, err = JsonFormatter(data)
		w.Write(jsonBytes)
	}

	w.Flush()
	return b.Bytes(), err
}

func tabulateOrganizationCollection(w *tabwriter.Writer, data *sdk.PaginatedOrganizationCollection) {
	fmt.Fprintf(w, "%v\t %v\t %v\n", "OrganizationID", "Name", "Timezone")

	if data.Embedded != nil {
		for _, item := range data.Embedded.Item {
			fmt.Fprintf(w, "%v\t %v\t %v\n", item.GetId(), item.GetName(), item.GetTimezone())
		}
	}
}

func tabulateOrganizationItem(w *tabwriter.Writer, item *sdk.OrganizationItem) {
	fmt.Fprintf(w, "%v\t %v\n", "OrganizationID", item.GetId())
	fmt.Fprintf(w, "%v\t %v\n", "Name", item.GetName())
	fmt.Fprintf(w, "%v\t %v\n", "Timezone", item.GetTimezone())
	fmt.Fprintf(w, "%v\t %v\n", "Projects", item.GetTotalProjects())
	fmt.Fprintf(w, "%v\t %v\n", "Clusters", item.GetAvailableClusters())
	fmt.Fprintf(w, "%v\t %v\n", "GitIntegrations", item.GetAvailableGitIntegration())
	fmt.Fprintf(w, "%v\t %v\n", "Registries", item.GetAvailableRegistries())
}

func tabulateProjectCollection(w *tabwriter.Writer, data *sdk.PaginatedProjectCollection) {
	fmt.Fprintf(w, "%v\t %v\t %v\t %v\n", "ProjectID", "OrganizationID", "Name", "Environments")

	if data.Embedded != nil {
		for _, item := range data.Embedded.Item {
			fmt.Fprintf(w, "%v\t %v\t %v\t %v\n", item.GetId(), item.GetOrganization(), item.GetName(), item.GetTotalEnvironments())
		}
	}
}

func tabulateProjectItem(w *tabwriter.Writer, item *sdk.ProjectItem) {
	fmt.Fprintf(w, "%v\t %v\n", "ProjectID", item.GetId())
	fmt.Fprintf(w, "%v\t %v\n", "OrganizationID", item.GetOrganization())
	fmt.Fprintf(w, "%v\t %v\n", "Name", item.GetName())
	fmt.Fprintf(w, "%v\t %v\n", "Environments", item.GetTotalEnvironments())
}

func tabulateEnvironmentCollection(w *tabwriter.Writer, data *sdk.PaginatedEnvironmentCollection) {
	fmt.Fprintf(w, "%v\t %v\t %v\t %v\t %v\t %v\n", "EnvironmentID", "ProjectID", "Name", "Namespace", "Type", "OperationStatus")

	if data.Embedded != nil {
		for _, item := range data.Embedded.Item {
			fmt.Fprintf(w, "%v\t %v\t %v\t %v\t %v\t %v\n", item.GetId(), item.GetProject(), item.GetName(), item.GetNamespace(), item.GetType(), item.GetOperationStatus())
		}
	}
}

func tabulateEnvironmentItem(w *tabwriter.Writer, item *sdk.EnvironmentItem) {
	fmt.Fprintf(w, "%v\t %v\n", "EnvironmentID", item.GetId())
	fmt.Fprintf(w, "%v\t %v\n", "ProjectID", item.GetProject())
	fmt.Fprintf(w, "%v\t %v\n", "Name", item.GetName())
	fmt.Fprintf(w, "%v\t %v\n", "Namespace", item.GetNamespace())
	fmt.Fprintf(w, "%v\t %v\n", "Type", item.GetType())
	fmt.Fprintf(w, "%v\t %v\n", "Components", item.GetTotalComponents())
	fmt.Fprintf(w, "%v\t %v\n", "OperationStatus", item.GetOperationStatus())
}

func tabulateComponentCollection(w *tabwriter.Writer, data *sdk.PaginatedComponentCollection) {
	fmt.Fprintf(w, "%v\t %v\t %v\t %v\t %v\n", "ComponentID", "EnvironmentID", "Name", "OperationStatus", "ClusterStatus")

	if data.Embedded != nil {
		for _, item := range data.Embedded.Item {
			fmt.Fprintf(w, "%v\t %v\t %v\t %v\t %v\n", item.GetId(), item.GetEnvironment(), item.GetName(), item.GetOperationStatus(), item.GetClusterStatus())
		}
	}
}

func tabulateComponentItem(w *tabwriter.Writer, item *sdk.ComponentItem) {
	fmt.Fprintf(w, "%v\t %v\n", "ComponentID", item.GetId())
	fmt.Fprintf(w, "%v\t %v\n", "EnvironmentID", item.GetEnvironment())
	fmt.Fprintf(w, "%v\t %v\n", "Name", item.GetName())
	fmt.Fprintf(w, "%v\t %v\n", "OperationStatus", item.GetOperationStatus())
	fmt.Fprintf(w, "%v\t %v\n", "ClusterStatus", item.GetClusterStatus())

	for index, url := range item.GetPublicURLs() {
		if index == 0 {
			fmt.Fprintf(w, "%v\t %v\n", "Public URL", url)
		} else {
			fmt.Fprintf(w, "\t %v\n", url)
		}
	}
}

func tabulateEventCollection(w *tabwriter.Writer, data *sdk.PaginatedEventCollection) {
	fmt.Fprintf(w, "%v\t %v\t %v\t %v\t %v\n", "EventID", "EnvironmentID", "OrganizationID", "Type", "Status")

	if data.Embedded != nil {
		for _, item := range data.Embedded.Item {
			fmt.Fprintf(w, "%v\t %v\t %v\t %v\t %v\n", item.GetId(), item.GetEnvironment(), item.GetOrganization(), item.GetType(), item.GetStatus())
		}
	}
}

func tabulateEventItem(w *tabwriter.Writer, item *sdk.EventItem) {
	fmt.Fprintf(w, "%v\t %v\n", "EventID", item.GetId())
	fmt.Fprintf(w, "%v\t %v\n", "EnvironmentID", item.GetEnvironment())
	fmt.Fprintf(w, "%v\t %v\n", "OrganizationID", item.GetOrganization())
	fmt.Fprintf(w, "%v\t %v\n", "Status", item.GetStatus())
	fmt.Fprintf(w, "%v\t %v\n", "Type", item.GetType())
	fmt.Fprintf(w, "%v\t %v\n", "CreatedAt", item.GetCreatedAt())
	fmt.Fprintf(w, "%v\t %v\n", "UpdatedAt", item.GetUpdatedAt())
}

func tabulateGeneric(w *tabwriter.Writer, item *sdk.ProblemGeneric) {
	fmt.Fprintf(w, "%v\n", "ERROR")
	fmt.Fprintf(w, "%v\t %v\n", "Title", item.GetTitle())
	fmt.Fprintf(w, "%v\t %v\n", "Detail", item.GetDetail())
}

func tabulatePagination(w *tabwriter.Writer, page int32, perPage int32, total int32) {
	var pages = total/perPage + 1
	if page > pages {
		fmt.Fprint(w, "\nPage does not exist")
	} else {
		fmt.Fprintf(w, "\nPage %d/%d with %d results", page, pages, total)
	}
}
