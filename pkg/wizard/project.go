package wizard

import (
	"fmt"

	"bunnyshell.com/sdk"
)

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
		return nil, ErrEmptyListing
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
	ctx, cancel := w.getContext()
	defer cancel()

	request := w.client.ProjectApi.ProjectList(ctx)

	if page > 1 {
		request = request.Page(page)
	}

	if w.profile.Context.Organization != "" {
		request = request.Organization(w.profile.Context.Organization)
	}

	paginated, _, err := request.Execute()

	return paginated, err
}
