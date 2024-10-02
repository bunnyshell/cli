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
	case *sdk.PaginatedEnvironItemCollection:
		tabulateEnvironCollection(writer, dataType)
	case *sdk.PaginatedProjectVariableCollection:
		tabulateProjectVariableCollection(writer, dataType)
	case *sdk.PaginatedKubernetesIntegrationCollection:
		tabulateKubernetesCollection(writer, dataType)
	case *sdk.PaginatedPipelineCollection:
		tabulatePipelineCollection(writer, dataType)
	case *sdk.PaginatedComponentGitCollection:
		tabulateComponentGitCollection(writer, dataType)
	case []sdk.ComponentGitCollection:
		tabulateComponentGitList(writer, dataType)
	case *sdk.PaginatedTemplateCollection:
		tabulateTemplateCollection(writer, dataType)
	case *sdk.PaginatedTemplatesRepositoryCollection:
		tabulateTemplatesRepositoryCollection(writer, dataType)
	case *sdk.PaginatedRegistryIntegrationCollection:
		tabulateRegistryIntegrationsCollection(writer, dataType)
	case *sdk.PaginatedServiceComponentVariableCollection:
		tabulateServiceComponentVariableCollection(writer, dataType)
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
	case *sdk.EnvironItemItem:
		tabulateEnvironItem(writer, dataType)
	case *sdk.ProjectVariableItem:
		tabulateProjectVariableItem(writer, dataType)
	case *sdk.ServiceComponentVariableItem:
		tabulateServiceComponentVariableItem(writer, dataType)
	case *sdk.KubernetesIntegrationItem:
		tabulateKubernetesItem(writer, dataType)
	case *sdk.RegistryIntegrationItem:
		tabulateRegistryIntegrationItem(writer, dataType)
	case *sdk.PipelineItem:
		tabulatePipelineItem(writer, dataType)
	case *sdk.ComponentGitItem:
		tabulateComponentGitItem(writer, dataType)
	case *sdk.SecretDecryptedItem:
		tabulateSecretDecryptedItem(writer, dataType)
	case *sdk.SecretEncryptedItem:
		tabulateSecretEncryptedItem(writer, dataType)
	case *sdk.TemplateItem:
		tabulateTemplateItem(writer, dataType)
	case *sdk.TemplatesRepositoryItem:
		tabulateTemplatesRepositoryItem(writer, dataType)
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

	first := true
	for key, value := range item.GetLabels() {
		if first {
			fmt.Fprintf(w, "%v\t %v\t %v\n", "Labels", key, value)

			first = false
		} else {
			fmt.Fprintf(w, "\t %v\t %v\n", key, value)
		}
	}

	if buildSettings, ok := item.GetBuildSettingsOk(); ok {
		tabulateBuildSettings(w, buildSettings)
	}
}

func tabulateBuildSettings(w *tabwriter.Writer, item *sdk.BuildSettingsItem) {
	buildCluster := ""
	if useManagedCluster, ok := item.GetUseManagedClusterOk(); ok && !*useManagedCluster {
		if k8sCluster, ok := item.GetKubernetesIntegrationOk(); ok && k8sCluster != nil {
			buildCluster = *k8sCluster
		} else {
			buildCluster = "Project cluster"
		}
	}

	if buildCluster != "" {
		fmt.Fprintf(w, "%v\t %v\n", "Build Cluster", buildCluster)
	}

	registryIntegration := ""
	if useManagedRegistry, ok := item.GetUseManagedRegistryOk(); ok && !*useManagedRegistry {
		if registry, ok := item.GetRegistryIntegrationOk(); ok && registry != nil {
			registryIntegration = *registry
		} else {
			registryIntegration = "Project registry"
		}
	}

	if registryIntegration != "" {
		fmt.Fprintf(w, "%v\t %v\n", "Build Registry", registryIntegration)
	}
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

	first := true
	for key, value := range item.GetLabels() {
		if first {
			fmt.Fprintf(w, "%v\t %v\t %v\n", "Labels", key, value)

			first = false
		} else {
			fmt.Fprintf(w, "\t %v\t %v\n", key, value)
		}
	}

	if buildSettings, ok := item.GetBuildSettingsOk(); ok {
		if buildSettings == nil {
			fmt.Fprintf(w, "%v\t %v\n", "Build Cluster", "Project cluster")
			fmt.Fprintf(w, "%v\t %v\n", "Build Registry", "Project registry")
		} else {
			tabulateBuildSettings(w, buildSettings)
		}
	}

	if item.GetType() == "primary" {
		fmt.Fprintf(w, "\n%s\n", "Primary Environment Settings")
		fmt.Fprintf(w, "%v\t %v\n", "Create Ephemeral On PR", item.GetHasEphemeralCreateOnPr())
		fmt.Fprintf(w, "%v\t %v\n", "Destroy Ephemeral On PR Close", item.GetHasEphemeralDestroyOnPrClose())
		fmt.Fprintf(w, "%v\t %v\n", "Auto Deploy Ephemeral", item.GetHasEphemeralAutoDeploy())
		fmt.Fprintf(w, "%v\t %v\n", "Termination Protection", item.GetHasTerminationProtection())

		if item.GetHasEphemeralBranchWhitelist() {
			fmt.Fprintf(w, "%v\t %v\n", "Ephemeral Branch Whitelist", item.GetEphemeralBranchWhitelistRegex())
		}
	}
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

func tabulateServiceComponentVariableCollection(w *tabwriter.Writer, data *sdk.PaginatedServiceComponentVariableCollection) {
	fmt.Fprintf(w, "%v\t %v\t %v\t %v\t %v\n", "ComponentVarID", "ComponentID", "EnvironmentId", "ProjectID", "Name")

	if data.Embedded != nil {
		for _, item := range data.Embedded.Item {
			fmt.Fprintf(w, "%v\t %v\t %v\t %v\t %v\n", item.GetId(), item.GetServiceComponent(), item.GetEnvironment(), item.GetProject(), item.GetName())
		}
	}
}

func tabulateServiceComponentVariableItem(w *tabwriter.Writer, item *sdk.ServiceComponentVariableItem) {
	fmt.Fprintf(w, "%v\t %v\n", "ComponentVariableID", item.GetId())
	fmt.Fprintf(w, "%v\t %v\n", "ComponentID", item.GetServiceComponent())
	fmt.Fprintf(w, "%v\t %v\n", "EnvironmentID", item.GetEnvironment())
	fmt.Fprintf(w, "%v\t %v\n", "ProjectID", item.GetProject())
	fmt.Fprintf(w, "%v\t %v\n", "Name", item.GetName())
	fmt.Fprintf(w, "%v\t %v\n", "Value", item.GetValue())
	fmt.Fprintf(w, "%v\t %v\n", "Secret", item.GetSecret())
}

func tabulateEnvironmentVariableCollection(w *tabwriter.Writer, data *sdk.PaginatedEnvironmentVariableCollection) {
	fmt.Fprintf(w, "%v\t %v\t %v\t %v\n", "EnvVarID", "EnvironmentID", "OrganizationID", "Name")

	if data.Embedded != nil {
		for _, item := range data.Embedded.Item {
			fmt.Fprintf(w, "%v\t %v\t %v\t %v\n", item.GetId(), item.GetEnvironment(), item.GetOrganization(), item.GetName())
		}
	}
}

func tabulateProjectVariableCollection(w *tabwriter.Writer, data *sdk.PaginatedProjectVariableCollection) {
	fmt.Fprintf(w, "%v\t %v\t %v\t %v\n", "ProjectVarID", "ProjectID", "OrganizationID", "Name")

	if data.Embedded != nil {
		for _, item := range data.Embedded.Item {
			fmt.Fprintf(w, "%v\t %v\t %v\t %v\n", item.GetId(), item.GetProject(), item.GetOrganization(), item.GetName())
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
	fmt.Fprintf(w, "%v\t %v\n", "EnvironmentVariableID", item.GetId())
	fmt.Fprintf(w, "%v\t %v\n", "EnvironmentID", item.GetEnvironment())
	fmt.Fprintf(w, "%v\t %v\n", "OrganizationID", item.GetOrganization())
	fmt.Fprintf(w, "%v\t %v\n", "Name", item.GetName())
	fmt.Fprintf(w, "%v\t %v\n", "Value", item.GetValue())
	fmt.Fprintf(w, "%v\t %v\n", "Secret", item.GetSecret())
}

func tabulateProjectVariableItem(w *tabwriter.Writer, item *sdk.ProjectVariableItem) {
	fmt.Fprintf(w, "%v\t %v\n", "ProjectVariableID", item.GetId())
	fmt.Fprintf(w, "%v\t %v\n", "ProjectID", item.GetProject())
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

	if len(item.Violations) == 0 {
		return
	}

	fmt.Fprintf(w, "\n%v\n", "VIOLATIONS")
	fmt.Fprintf(w, "%v\t %v\n", "Property", "Message")

	for _, violation := range item.Violations {
		fmt.Fprintf(w, "%v\t %v\n", violation.GetPropertyPath(), violation.GetMessage())
	}
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
