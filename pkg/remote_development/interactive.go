package remote_development

import (
	"fmt"
	"strings"

	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/util"
	bunnysdk "bunnyshell.com/sdk"
)

var (
	ErrNoEnvironments        = fmt.Errorf("no environments available")
	ErrNoOrganizations       = fmt.Errorf("no organizations available")
	ErrNoComponents          = fmt.Errorf("no components available")
	ErrNoComponentResourcess = fmt.Errorf("no component resourcess available")
	ErrNoProjects            = fmt.Errorf("no projects available")
)

func (r *RemoteDevelopment) SelectOrganization() error {
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

func (r *RemoteDevelopment) SelectProject() error {
	resp, _, err := lib.GetProjects(r.organization.GetId())
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

func (r *RemoteDevelopment) SelectEnvironment() error {
	resp, _, err := lib.GetEnvironments(r.project.GetId())
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

func (r *RemoteDevelopment) SelectComponent() error {
	resp, _, err := lib.GetComponents(r.environment.GetId(), "running")

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

func (r *RemoteDevelopment) SelectComponentResource() error {
	resources, _, err := lib.GetComponentResources(r.component.GetId())
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
