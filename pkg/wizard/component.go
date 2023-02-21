package wizard

import (
	"fmt"

	"bunnyshell.com/cli/pkg/api/component"
	"bunnyshell.com/sdk"
)

func (w *Wizard) SelectComponent() (*sdk.ComponentCollection, error) {
	return w.selectComponent(1)
}

func (w *Wizard) selectComponent(page int32) (*sdk.ComponentCollection, error) {
	model, err := w.getComponents(page)
	if err != nil {
		return nil, err
	}

	embedded, ok := model.GetEmbeddedOk()
	if !ok {
		return nil, ErrEmptyListing
	}

	collectionItems := embedded.GetItem()

	items := []string{}
	for _, item := range collectionItems {
		items = append(items, fmt.Sprintf("%s (%s)", item.GetName(), item.GetId()))
	}

	currentPage, totalPages := getPaginationInfo(model)

	index, newPage, err := chooseOrNavigate("Select Component", items, currentPage, totalPages)
	if err != nil {
		return nil, err
	}

	if newPage != nil {
		return w.selectComponent(*newPage)
	}

	if index != nil {
		return &collectionItems[*index], nil
	}

	panic("Something went wrong...")
}

func (w *Wizard) getComponents(page int32) (*sdk.PaginatedComponentCollection, error) {
	listOptions := component.NewListOptions()
	listOptions.Page = page
	listOptions.Profile = w.profile

	listOptions.Organization = w.profile.Context.Organization
	listOptions.Project = w.profile.Context.Project
	listOptions.Environment = w.profile.Context.Environment

	return component.List(listOptions)
}
