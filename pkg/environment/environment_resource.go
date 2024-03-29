package environment

import (
	"fmt"
	"strings"

	"bunnyshell.com/cli/pkg/api/component"
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
	allowedResouceTypes      = []ResourceType{Deployment, StatefulSet, DaemonSet}
	ErrNoEnvironments        = fmt.Errorf("no environments available")
	ErrNoOrganizations       = fmt.Errorf("no organizations available")
	ErrNoComponents          = fmt.Errorf("no components available")
	ErrNoComponentResourcess = fmt.Errorf("no component resourcess available")
	ErrNoProjects            = fmt.Errorf("no projects available")
)

type EnvironmentResource struct {
	Organization *bunnysdk.OrganizationItem
	Project      *bunnysdk.ProjectItem
	Environment  *bunnysdk.EnvironmentItem

	Component         *bunnysdk.ComponentItem
	ComponentResource *bunnysdk.ComponentResourceItem
}

func NewEnvironmentResource() *EnvironmentResource {
	return &EnvironmentResource{}
}

func (r *EnvironmentResource) WithComponent(component *bunnysdk.ComponentItem) *EnvironmentResource {
	if r.Environment == nil {
		panic(fmt.Errorf("you have to select an environment before selecting a component"))
	}

	if r.Environment.GetId() != component.GetEnvironment() {
		panic(fmt.Errorf(
			"component \"%s\" is not part of environment \"%s\"",
			component.GetName(),
			r.Environment.GetName(),
		))
	}

	r.Component = component

	return r
}

func (r *EnvironmentResource) WithComponentResource(component *bunnysdk.ComponentResourceItem) *EnvironmentResource {
	r.ComponentResource = component

	return r
}

func (r *EnvironmentResource) WithResourcePath(resourcePath string) *EnvironmentResource {
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

	resources, err := component.Resources(component.NewResourceOptions(r.Component.GetId()))
	if err != nil {
		panic(fmt.Errorf(
			"failed fetching resources for component \"%s\"",
			r.Component.GetId(),
		))
	}

	for _, resourceItem := range resources {
		if resourceItem.GetNamespace() == namespace && strings.ToLower(resourceItem.GetKind()) == resourceType && resourceItem.GetName() == resourceName {
			r.WithComponentResource(&resourceItem)

			break
		}
	}

	if r.ComponentResource == nil {
		panic(fmt.Errorf(
			"the component does not contain the \"%s\" resource",
			resourcePath,
		))
	}

	return r
}

func (r *EnvironmentResource) WithOrganization(organization *bunnysdk.OrganizationItem) *EnvironmentResource {
	r.Organization = organization

	return r
}

func (r *EnvironmentResource) WithProject(project *bunnysdk.ProjectItem) *EnvironmentResource {
	if r.Organization != nil && r.Organization.GetId() != project.GetOrganization() {
		panic(fmt.Errorf(
			"project \"%s\" is not part of organization \"%s\"",
			project.GetName(),
			r.Organization.GetName(),
		))
	}

	r.Project = project

	return r
}

func (r *EnvironmentResource) WithEnvironment(environment *bunnysdk.EnvironmentItem) *EnvironmentResource {
	if r.Project != nil && r.Project.GetId() != environment.GetProject() {
		panic(fmt.Errorf(
			"environment \"%s\" is not part of project \"%s\"",
			environment.GetName(),
			r.Project.GetName(),
		))
	}

	r.Environment = environment

	return r
}
