package environment

import (
	"fmt"
	"strings"

	"bunnyshell.com/cli/pkg/api/component"
	"bunnyshell.com/cli/pkg/api/environment"
	"bunnyshell.com/cli/pkg/api/organization"
	"bunnyshell.com/cli/pkg/api/project"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/interactive"
	bunnysdk "bunnyshell.com/sdk"
)

func NewFromWizard(profileContext *config.Context, resourcePath string) (*EnvironmentResource, error) {
	environmentResource, err := getEnvironmentResource(profileContext)
	if err != nil {
		return nil, err
	}

	if resourcePath != "" {
		return environmentResource.WithResourcePath(resourcePath), nil
	}

	if config.GetSettings().NonInteractive {
		return nil, interactive.ErrNonInteractive
	}

	if err = environmentResource.SelectComponentResource(); err != nil {
		return nil, err
	}

	return environmentResource, nil
}

func getEnvironmentResource(profileContext *config.Context) (*EnvironmentResource, error) {
	// wizard
	if profileContext.ServiceComponent != "" {
		componentItem, err := component.Get(component.NewItemOptions(profileContext.ServiceComponent))
		if err != nil {
			return nil, err
		}

		environmentItem, err := environment.Get(environment.NewItemOptions(componentItem.GetEnvironment()))
		if err != nil {
			return nil, err
		}

		environmentResource := NewEnvironmentResource().WithEnvironment(environmentItem).WithComponent(componentItem)

		return environmentResource, nil
	}

	if config.GetSettings().NonInteractive {
		return nil, interactive.ErrNonInteractive
	}

	return askEnvironmentResource(profileContext)
}

func askEnvironmentResource(profileContext *config.Context) (*EnvironmentResource, error) {
	environmentResource := NewEnvironmentResource()

	if profileContext.Organization != "" {
		itemOptions := organization.NewItemOptions(profileContext.Organization)

		organizationItem, err := organization.Get(itemOptions)
		if err != nil {
			return nil, err
		}

		environmentResource.WithOrganization(organizationItem)
	} else if err := environmentResource.SelectOrganization(); err != nil {
		return nil, err
	}

	if profileContext.Project != "" {
		itemOptions := project.NewItemOptions(profileContext.Project)

		projectItem, err := project.Get(itemOptions)
		if err != nil {
			return nil, err
		}

		environmentResource.WithProject(projectItem)
	} else if err := environmentResource.SelectProject(); err != nil {
		return nil, err
	}

	if profileContext.Environment != "" {
		itemOptions := environment.NewItemOptions(profileContext.Environment)

		environmentItem, err := environment.Get(itemOptions)
		if err != nil {
			return nil, err
		}

		environmentResource.WithEnvironment(environmentItem)
	} else if err := environmentResource.SelectEnvironment(); err != nil {
		return nil, err
	}

	if err := environmentResource.SelectComponent(); err != nil {
		return nil, err
	}

	return environmentResource, nil
}

func (r *EnvironmentResource) SelectOrganization() error {
	model, err := organization.List(nil)
	if err != nil {
		return err
	}

	if model.Embedded == nil {
		return ErrNoOrganizations
	}

	items := []string{}
	for _, item := range model.Embedded.GetItem() {
		items = append(items, item.GetName())
	}

	index, _, err := interactive.Choose("Select organization", items)
	if err != nil {
		return err
	}

	itemOptions := organization.NewItemOptions(model.Embedded.GetItem()[index].GetId())

	organizationItem, err := organization.Get(itemOptions)
	if err != nil {
		return err
	}

	r.WithOrganization(organizationItem)

	return nil
}

func (r *EnvironmentResource) SelectProject() error {
	listOptions := project.NewListOptions()
	listOptions.Organization = r.Organization.GetId()

	resp, err := project.List(listOptions)
	if err != nil {
		return err
	}

	if resp.Embedded == nil {
		return ErrNoProjects
	}

	items := []string{}
	for _, item := range resp.Embedded.GetItem() {
		items = append(items, item.GetName())
	}

	index, _, err := interactive.Choose("Select project", items)
	if err != nil {
		return err
	}

	itemOptions := project.NewItemOptions(resp.Embedded.GetItem()[index].GetId())

	projectItem, err := project.Get(itemOptions)
	if err != nil {
		return err
	}

	r.WithProject(projectItem)

	return nil
}

func (r *EnvironmentResource) SelectEnvironment() error {
	listOptions := environment.NewListOptions()
	listOptions.Project = r.Project.GetId()

	resp, err := environment.List(listOptions)
	if err != nil {
		return err
	}

	if resp.Embedded == nil {
		return ErrNoEnvironments
	}

	items := []string{}
	for _, item := range resp.Embedded.GetItem() {
		items = append(items, item.GetName())
	}

	index, _, err := interactive.Choose("Select environment", items)
	if err != nil {
		return err
	}

	itemOptions := environment.NewItemOptions(resp.Embedded.GetItem()[index].GetId())

	environmentItem, err := environment.Get(itemOptions)
	if err != nil {
		return err
	}

	r.WithEnvironment(environmentItem)

	return nil
}

func (r *EnvironmentResource) SelectComponent() error {
	listOptions := component.NewListOptions()
	listOptions.Environment = r.Environment.GetId()
	listOptions.OperationStatus = "running"

	resp, err := component.List(listOptions)
	if err != nil {
		return err
	}

	if resp.Embedded == nil {
		return ErrNoComponents
	}

	components := resp.Embedded.GetItem()

	items := []string{}
	for _, item := range components {
		items = append(items, item.GetName())
	}

	index, _, err := interactive.Choose("Select component", items)
	if err != nil {
		return err
	}

	itemOptions := component.NewItemOptions(components[index].GetId())

	componentItem, err := component.Get(itemOptions)
	if err != nil {
		return err
	}

	r.WithComponent(componentItem)

	return nil
}

func (r *EnvironmentResource) SelectComponentResource() error {
	resources, err := component.Resources(component.NewResourceOptions(r.Component.GetId()))
	if err != nil {
		return err
	}

	allowedResouceTypesSet := map[string]bool{}
	for _, v := range allowedResouceTypes {
		allowedResouceTypesSet[strings.ToLower(string(v))] = true
	}

	allowedResources := []bunnysdk.ComponentResourceItem{}
	selectItems := []string{}

	for _, item := range resources {
		itemKind := strings.ToLower(item.GetKind())
		if _, present := allowedResouceTypesSet[itemKind]; !present {
			continue
		}

		allowedResources = append(allowedResources, item)
		selectItems = append(selectItems, fmt.Sprintf("%s / %s / %s", item.GetNamespace(), item.GetKind(), item.GetName()))
	}

	if len(allowedResources) == 0 {
		return ErrNoComponentResourcess
	}

	if len(allowedResources) == 1 {
		r.WithComponentResource(&allowedResources[0])

		return nil
	}

	index, _, err := interactive.Choose("Select resource", selectItems)
	if err != nil {
		return err
	}

	r.WithComponentResource(&allowedResources[index])

	return nil
}
