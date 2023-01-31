package wizard

import (
	"fmt"

	"bunnyshell.com/cli/pkg/api/organization"
	"bunnyshell.com/sdk"
)

func (w *Wizard) SelectOrganization() (*sdk.OrganizationCollection, error) {
	return w.selectOrganization(1)
}

func (w *Wizard) selectOrganization(page int32) (*sdk.OrganizationCollection, error) {
	model, err := w.getOrganizations(page)
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

	index, newPage, err := chooseOrNavigate("Select Organization", items, currentPage, totalPages)
	if err != nil {
		return nil, err
	}

	if newPage != nil {
		return w.selectOrganization(*newPage)
	}

	if index != nil {
		return &collectionItems[*index], nil
	}

	panic("Something went wrong...")
}

func (w *Wizard) getOrganizations(page int32) (*sdk.PaginatedOrganizationCollection, error) {
	listOptions := organization.NewListOptions()
	listOptions.Page = page
	listOptions.Profile = w.profile

	return organization.List(listOptions)
}
