package wizard

import (
	"fmt"

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
	ctx, cancel := w.getContext()
	defer cancel()

	request := w.client.ComponentApi.ComponentList(ctx)

	if page > 1 {
		request = request.Page(page)
	}

	if w.profile.Context.Environment != "" {
		request = request.Environment(w.profile.Context.Environment)
	}

	if w.profile.Context.Project != "" {
		request = request.Project(w.profile.Context.Project)
	}

	if w.profile.Context.Organization != "" {
		request = request.Organization(w.profile.Context.Organization)
	}

	paginated, _, err := request.Execute()

	return paginated, err
}
