package wizard

import (
	"bunnyshell.com/cli/pkg/interactive"
)

// @see lib/pagination.go
type Pagination interface {
	GetPage() int32

	GetItemsPerPage() int32

	GetTotalItems() int32
}

func getPaginationInfo(pagination Pagination) (int32, int32) {
	page := pagination.GetPage()
	pages := 1 + (pagination.GetTotalItems()-1)/pagination.GetItemsPerPage()

	return page, pages
}

func addPagination(items []string, page int32, pages int32) []string {
	if page > pages {
		return items
	}

	if pages == 1 {
		return items
	}

	items = append(items, paginationOptions[selectPage])

	return items
}

func chooseOrNavigate(question string, items []string, page int32, pages int32) (*int, *int32, error) {
	selector := addPagination(items, page, pages)

	index, answer, err := interactive.ChooseWithSize(paginationDisplaySize, question, selector)
	if err != nil {
		return nil, nil, err
	}

	if index >= len(items) {
		if answer == paginationOptions[selectPage] {
			page, err = interactive.AskInt32("Select a page", interactive.AssertBetween(1, pages))

			return nil, &page, err
		}

		return nil, nil, interactive.ErrInvalidValue
	}

	return &index, nil, nil
}
