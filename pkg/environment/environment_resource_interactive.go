package environment

import (
	"fmt"
	"strings"

	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/util"
	bunnysdk "bunnyshell.com/sdk"
)

func NewFromWizard(profileContext *lib.Context, resourcePath string) (*EnvironmentResource, error) {
	environmentResource := NewEnvironmentResource()

	// wizard
	if profileContext.ServiceComponent != "" {
		componentItem, _, err := lib.GetComponent(profileContext.ServiceComponent)
		if err != nil {
			return nil, err
		}

		environmentItem, _, err := lib.GetEnvironment(componentItem.GetEnvironment())
		if err != nil {
			return nil, err
		}

		environmentResource.WithEnvironment(environmentItem).WithComponent(componentItem)
	} else {
		if profileContext.Organization != "" {
			organizationItem, _, err := lib.GetOrganization(profileContext.Organization)
			if err != nil {
				return nil, err
			}

			environmentResource.WithOrganization(organizationItem)
		} else if err := environmentResource.SelectOrganization(); err != nil {
			return nil, err
		}

		if profileContext.Project != "" {
			projectItem, _, err := lib.GetProject(profileContext.Project)
			if err != nil {
				return nil, err
			}

			environmentResource.WithProject(projectItem)
		} else if err := environmentResource.SelectProject(); err != nil {
			return nil, err
		}

		if profileContext.Environment != "" {
			environmentItem, _, err := lib.GetEnvironment(profileContext.Environment)
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
	}

	if resourcePath != "" {
		environmentResource.WithResourcePath(resourcePath)
	} else {
		if err := environmentResource.SelectComponentResource(); err != nil {
			return nil, err
		}
	}

	return environmentResource, nil
}

func (r *EnvironmentResource) SelectOrganization() error {
	resp, _, err := lib.GetOrganizations()
	if err != nil {
		return err
	}

	if resp.Embedded == nil {
		return ErrNoOrganizations
	}

	items := []string{}
	for _, item := range resp.Embedded.GetItem() {
		items = append(items, item.GetName())
	}
	index, _, err := util.Choose("Select organization", items)
	if err != nil {
		return err
	}

	organizationItem, _, err := lib.GetOrganization(resp.Embedded.GetItem()[index].GetId())
	if err != nil {
		return err
	}

	r.WithOrganization(organizationItem)
	return nil
}

func (r *EnvironmentResource) SelectProject() error {
	resp, _, err := lib.GetProjects(r.Organization.GetId())
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
	index, _, err := util.Choose("Select project", items)
	if err != nil {
		return err
	}

	projectItem, _, err := lib.GetProject(resp.Embedded.GetItem()[index].GetId())
	if err != nil {
		return err
	}

	r.WithProject(projectItem)
	return nil
}

func (r *EnvironmentResource) SelectEnvironment() error {
	resp, _, err := lib.GetEnvironments(r.Project.GetId())
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
	index, _, err := util.Choose("Select environment", items)
	if err != nil {
		return err
	}

	environmentItem, _, err := lib.GetEnvironment(resp.Embedded.GetItem()[index].GetId())
	if err != nil {
		return err
	}

	r.WithEnvironment(environmentItem)
	return nil
}

func (r *EnvironmentResource) SelectComponent() error {
	resp, _, err := lib.GetComponents(r.Environment.GetId(), "running")

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
	index, _, err := util.Choose("Select component", items)
	if err != nil {
		return err
	}

	componentItem, _, err := lib.GetComponent(components[index].GetId())
	if err != nil {
		return err
	}

	r.WithComponent(componentItem)
	return nil
}

func (r *EnvironmentResource) SelectComponentResource() error {
	resources, _, err := lib.GetComponentResources(r.Component.GetId())
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

	index, _, err := util.Choose("Select resource", selectItems)
	if err != nil {
		return err
	}

	r.WithComponentResource(&allowedResources[index])
	return nil
}
