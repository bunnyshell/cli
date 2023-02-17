package formatter

import (
	"bytes"
	"fmt"
	"text/tabwriter"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/sdk"
)

func stylish(data interface{}) ([]byte, error) {
	var (
		buffer bytes.Buffer
		err    error
	)

	writer := tabwriter.NewWriter(&buffer, 1, 1, 1, ' ', tabwriter.Debug)

	switch dataType := data.(type) {
	case *sdk.PaginatedOrganizationCollection:
		tabulateOrganizationCollection(writer, dataType)
	case *sdk.PaginatedProjectCollection:
		tabulateProjectCollection(writer, dataType)
	case *sdk.PaginatedEnvironmentCollection:
		tabulateEnvironmentCollection(writer, dataType)
	case *sdk.PaginatedComponentCollection:
		tabulateComponentCollection(writer, dataType)
	case *sdk.PaginatedEventCollection:
		tabulateEventCollection(writer, dataType)
	case *sdk.PaginatedEnvironmentVariableCollection:
		tabulateEnvironmentVariableCollection(writer, dataType)
	case *sdk.PaginatedKubernetesIntegrationCollection:
		tabulateKubernetesCollection(writer, dataType)
	case *sdk.PaginatedPipelineCollection:
		tabulatePipelineCollection(writer, dataType)
	case *sdk.PaginatedComponentGitCollection:
		tabulateComponentGitCollection(writer, dataType)
	case []sdk.ComponentEndpointCollection:
		tabulateAggregateEndpoint(writer, dataType)
	case *sdk.OrganizationItem:
		tabulateOrganizationItem(writer, dataType)
	case *sdk.ProjectItem:
		tabulateProjectItem(writer, dataType)
	case *sdk.EnvironmentItem:
		tabulateEnvironmentItem(writer, dataType)
	case *sdk.ComponentItem:
		tabulateComponentItem(writer, dataType)
	case *sdk.EventItem:
		tabulateEventItem(writer, dataType)
	case *sdk.EnvironmentVariableItem:
		tabulateEnvironmentVariableItem(writer, dataType)
	case *sdk.KubernetesIntegrationItem:
		tabulateKubernetesItem(writer, dataType)
	case *sdk.PipelineItem:
		tabulatePipelineItem(writer, dataType)
	case *sdk.ComponentGitItem:
		tabulateComponentGitItem(writer, dataType)
	case *sdk.ProblemGeneric:
		tabulateGeneric(writer, dataType)
	case *api.Error:
		tabulateAPIError(writer, dataType)
	case api.Error:
		tabulateAPIError(writer, &dataType)
	case error:
		tabulateError(writer, dataType)
	default:
		err = writeJSON(writer, data)
	}

	writer.Flush()

	return buffer.Bytes(), err
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

func tabulateEnvironmentVariableCollection(w *tabwriter.Writer, data *sdk.PaginatedEnvironmentVariableCollection) {
	fmt.Fprintf(w, "%v\t %v\t %v\t %v\n", "EnvVarID", "EnvironmentID", "OrganizationID", "Name")

	if data.Embedded != nil {
		for _, item := range data.Embedded.Item {
			fmt.Fprintf(w, "%v\t %v\t %v\t %v\n", item.GetId(), item.GetEnvironment(), item.GetOrganization(), item.GetName())
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

func tabulateEnvironmentVariableItem(w *tabwriter.Writer, item *sdk.EnvironmentVariableItem) {
	fmt.Fprintf(w, "%v\t %v\n", "EnvironmentID", item.GetEnvironment())
	fmt.Fprintf(w, "%v\t %v\n", "OrganizationID", item.GetOrganization())
	fmt.Fprintf(w, "%v\t %v\n", "Name", item.GetName())
	fmt.Fprintf(w, "%v\t %v\n", "Value", item.GetValue())
	fmt.Fprintf(w, "%v\t %v\n", "Secret", item.GetSecret())
}

func tabulateGeneric(w *tabwriter.Writer, item *sdk.ProblemGeneric) {
	fmt.Fprintf(w, "%v\n", "ERROR")
	fmt.Fprintf(w, "%v\t %v\n", "Title", item.GetTitle())
	fmt.Fprintf(w, "%v\t %v\n", "Detail", item.GetDetail())
}

func tabulateAPIError(w *tabwriter.Writer, item *api.Error) {
	fmt.Fprintf(w, "%v\n", "ERROR")
	fmt.Fprintf(w, "%v\t %v\n", "Title", item.Title)
	fmt.Fprintf(w, "%v\t %v\n", "Detail", item.Detail)
}

func tabulateError(w *tabwriter.Writer, err error) {
	fmt.Fprintf(w, "\n%v\t %v\n", "ERROR", err.Error())
}

func writeJSON(writer *tabwriter.Writer, data any) error {
	fmt.Fprintf(writer, "JSON: ")

	jsonBytes, err := JSONFormatter(data)
	_, _ = writer.Write(jsonBytes)

	return err
}
