package wizard

import (
	"fmt"

	"bunnyshell.com/cli/pkg/api/environment"
	"bunnyshell.com/sdk"
)

func (w *Wizard) HasEnvironment() bool {
	return w.profile.Context.Environment != ""
}

func (w *Wizard) GetEnvironment() (*sdk.EnvironmentItem, error) {
	if !w.HasEnvironment() {
		itemCol, err := w.SelectEnvironment()
		if err != nil {
			return nil, err
		}

		w.profile.Context.Environment = itemCol.GetId()
	}

	item, err := environment.Get(environment.NewItemOptions(w.profile.Context.Environment))
	if err != nil {
		return nil, err
	}

	w.profile.Context.Project = item.GetProject()

	return item, nil
}

func (w *Wizard) SelectEnvironment() (*sdk.EnvironmentCollection, error) {
	return w.selectEnvironment(1)
}

func (w *Wizard) selectEnvironment(page int32) (*sdk.EnvironmentCollection, error) {
	model, err := w.getEnvironments(page)
	if err != nil {
		return nil, err
	}

	embedded, ok := model.GetEmbeddedOk()
	if !ok {
		return nil, fmt.Errorf("%s %w", "environments", ErrEmptyListing)
	}

	collectionItems := embedded.GetItem()

	items := []string{}
	for _, item := range collectionItems {
		items = append(items, fmt.Sprintf("%s (%s)", item.GetName(), item.GetId()))
	}

	currentPage, totalPages := getPaginationInfo(model)

	index, newPage, err := chooseOrNavigate("Select Environment", items, currentPage, totalPages)
	if err != nil {
		return nil, err
	}

	if newPage != nil {
		return w.selectEnvironment(*newPage)
	}

	if index != nil {
		return &collectionItems[*index], nil
	}

	panic("Something went wrong...")
}

func (w *Wizard) getEnvironments(page int32) (*sdk.PaginatedEnvironmentCollection, error) {
	listOptions := environment.NewListOptions()
	listOptions.Page = page
	listOptions.Organization = w.profile.Context.Organization
	listOptions.Project = w.profile.Context.Project
	listOptions.Profile = w.profile

	return environment.List(listOptions)
}
