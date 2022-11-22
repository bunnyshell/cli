package remote_development

import (
	"fmt"
	"strings"

	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/dev/pkg/remote"
	bunnysdk "bunnyshell.com/sdk"
)

// +enum
type ResourceType string

const (
	Deployment  ResourceType = "deployment"
	StatefulSet ResourceType = "statefulset"
	DaemonSet   ResourceType = "daemonset"
)

var (
	ErrNoOrganizationSelected = fmt.Errorf("you need to select an organization first")
	allowedResouceTypes       = []ResourceType{Deployment, StatefulSet, DaemonSet}
)

type RemoteDevelopment struct {
	remoteDev *remote.RemoteDevelopment

	organization      *bunnysdk.OrganizationItem
	project           *bunnysdk.ProjectItem
	environment       *bunnysdk.EnvironmentItem
	component         *bunnysdk.ComponentItem
	componentResource *bunnysdk.ComponentResourceItem

	environmentWorkspaceDir string

	kubeConfigPath string

	localSyncPath  string
	remoteSyncPath string

	portMappings []string
}

func NewRemoteDevelopment() *RemoteDevelopment {
	return &RemoteDevelopment{
		remoteDev: remote.NewRemoteDevelopment(),

		organization: nil,
		project:      nil,
		environment:  nil,
		component:    nil,
	}
}

func (r *RemoteDevelopment) WithOrganization(organization *bunnysdk.OrganizationItem) *RemoteDevelopment {
	r.organization = organization
	return r
}

func (r *RemoteDevelopment) WithProject(project *bunnysdk.ProjectItem) *RemoteDevelopment {
	if r.organization != nil && r.organization.GetId() != project.GetOrganization() {
		panic(fmt.Errorf(
			"project \"%s\" is not part of organization \"%s\"",
			project.GetName(),
			r.organization.GetName(),
		))
	}

	r.project = project
	return r
}

func (r *RemoteDevelopment) WithEnvironment(environment *bunnysdk.EnvironmentItem) *RemoteDevelopment {
	if r.project != nil && r.project.GetId() != environment.GetProject() {
		panic(fmt.Errorf(
			"environment \"%s\" is not part of project \"%s\"",
			environment.GetName(),
			r.project.GetName(),
		))
	}

	r.environment = environment
	return r
}

func (r *RemoteDevelopment) WithComponent(component *bunnysdk.ComponentItem) *RemoteDevelopment {
	if r.environment == nil {
		panic(fmt.Errorf("you have to select an environment before selecting a component"))
	}

	if r.environment.GetId() != component.GetEnvironment() {
		panic(fmt.Errorf(
			"component \"%s\" is not part of environment \"%s\"",
			component.GetName(),
			r.environment.GetName(),
		))
	}

	r.component = component
	return r
}

func (r *RemoteDevelopment) WithComponentResource(component *bunnysdk.ComponentResourceItem) *RemoteDevelopment {
	r.componentResource = component
	return r
}

func (r *RemoteDevelopment) WithEnvironmentWorkspaceDir(environmentWorkspaceDir string) *RemoteDevelopment {
	r.environmentWorkspaceDir = environmentWorkspaceDir
	return r
}

func (r *RemoteDevelopment) WithKubeConfigPath(kubeConfigPath string) *RemoteDevelopment {
	r.kubeConfigPath = kubeConfigPath
	return r
}

func (r *RemoteDevelopment) WithLocalSyncPath(localSyncPath string) *RemoteDevelopment {
	r.localSyncPath = localSyncPath
	return r
}

func (r *RemoteDevelopment) WithRemoteSyncPath(remoteSyncPath string) *RemoteDevelopment {
	r.remoteSyncPath = remoteSyncPath
	return r
}

func (r *RemoteDevelopment) WithResourcePath(resourcePath string) *RemoteDevelopment {
	resourceParts := strings.Split(resourcePath, "/")
	if len(resourceParts) != 3 {
		panic(fmt.Errorf(
			"the provided resource path \"%s\" is invalid",
			resourcePath,
		))
	}

	namespace := resourceParts[0]
	resourceType := strings.ToLower(resourceParts[1])
	resourceName := resourceParts[2]

	allowedResouceTypesSet := map[string]bool{}
	for _, v := range allowedResouceTypes {
		allowedResouceTypesSet[strings.ToLower(string(v))] = true
	}

	if _, present := allowedResouceTypesSet[resourceType]; !present {
		panic(fmt.Errorf(
			"the provided resource type \"%s\" is not valid for remote development",
			resourceType,
		))
	}

	resources, _, err := lib.GetComponentResources(r.component.GetId())
	if err != nil {
		panic(fmt.Errorf(
			"failed fetching resources for component \"%s\"",
			r.component.GetId(),
		))
	}

	for _, resourceItem := range resources {
		if resourceItem.GetNamespace() == namespace && strings.ToLower(resourceItem.GetKind()) == resourceType && resourceItem.GetName() == resourceName {
			r.WithComponentResource(&resourceItem)

			break
		}
	}
	if r.componentResource == nil {
		panic(fmt.Errorf(
			"the component does not contain the \"%s\" resource",
			resourcePath,
		))
	}

	return r
}

func (r *RemoteDevelopment) WithPortMappings(portMappings []string) *RemoteDevelopment {
	r.portMappings = portMappings
	return r
}
