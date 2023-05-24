package wizard

import (
	"fmt"

	"bunnyshell.com/cli/pkg/api/project"
	"bunnyshell.com/sdk"
)

func (w *Wizard) HasProject() bool {
	return w.profile.Context.Project != ""
}

func (w *Wizard) GetProject() (*sdk.ProjectItem, error) {
	if !w.HasProject() {
		itemCol, err := w.SelectProject()
		if err != nil {
			return nil, err
		}

		w.profile.Context.Project = itemCol.GetId()
	}

	item, err := project.Get(project.NewItemOptions(w.profile.Context.Project))
	if err != nil {
		return nil, err
	}

	w.profile.Context.Organization = item.GetOrganization()

	return item, nil
}

func (w *Wizard) SelectProject() (*sdk.ProjectCollection, error) {
	return w.selectProject(1)
}

func (w *Wizard) selectProject(page int32) (*sdk.ProjectCollection, error) {
	model, err := w.getProjects(page)
	if err != nil {
		return nil, err
	}

	embedded, ok := model.GetEmbeddedOk()
	if !ok {
		return nil, fmt.Errorf("%s %w", "projects", ErrEmptyListing)
	}

	collectionItems := embedded.GetItem()

	items := []string{}
	for _, item := range collectionItems {
		items = append(items, fmt.Sprintf("%s (%s)", item.GetName(), item.GetId()))
	}

	currentPage, totalPages := getPaginationInfo(model)

	index, newPage, err := chooseOrNavigate("Select Project", items, currentPage, totalPages)
	if err != nil {
		return nil, err
	}

	if newPage != nil {
		return w.selectProject(*newPage)
	}

	if index != nil {
		return &collectionItems[*index], nil
	}

	panic("Something went wrong...")
}

func (w *Wizard) getProjects(page int32) (*sdk.PaginatedProjectCollection, error) {
	listOptions := project.NewListOptions()
	listOptions.Page = page
	listOptions.Profile = w.profile

	listOptions.Organization = w.profile.Context.Organization

	return project.List(listOptions)
}
