package lib

import (
	"fmt"
	"net/http"

	"bunnyshell.com/cli/pkg/util"
	"github.com/spf13/cobra"
)

const (
	PAGINATION_QUIT = -1

	PAGINATION_OTHER = -2
)

type ModelWithPagination interface {
	GetPage() int32

	GetItemsPerPage() int32

	GetTotalItems() int32
}

type CollectionGenerator func(page int32) (ModelWithPagination, *http.Response, error)

func ShowCollection(cmd *cobra.Command, page int32, generator CollectionGenerator) error {
	for {
		model, resp, err := generator(page)
		if err = FormatRequestResult(cmd, model, resp, err); err != nil {
			return err
		}

		if CLIContext.OutputFormat != "stylish" {
			return nil
		}

		page, err := ProcessPagination(cmd, model)
		if err != nil {
			return err
		}

		if page == PAGINATION_QUIT {
			return nil
		}
	}
}

func ProcessPagination(cmd *cobra.Command, m ModelWithPagination) (int32, error) {
	page := m.GetPage()
	pages := 1 + (m.GetTotalItems()-1)/m.GetItemsPerPage()

	if page > pages {
		return PAGINATION_QUIT, nil
	}

	if pages == 1 {
		return PAGINATION_QUIT, nil
	}

	nav, pageNo := getPaginationOptions(page, pages)

	index, _, err := util.Choose("Navigate to a different page?", nav)
	if err != nil {
		return PAGINATION_QUIT, err
	}

	target := pageNo[index]
	if target == PAGINATION_OTHER {
		target, err = util.AskInt32("Go to page:", util.AssertBetween(1, pages))
		if err != nil {
			return PAGINATION_QUIT, err
		}
	}

	return target, nil
}

func getPaginationOptions(page int32, pages int32) ([]string, []int32) {
	nav := []string{}
	pageNo := []int32{}
	var firstPage int32 = 1
	var lastPage int32 = pages

	if page != firstPage {
		var prevPage int32 = page - 1

		if prevPage != firstPage {
			nav = append(nav, fmt.Sprintf("First (%d)", firstPage))
			pageNo = append(pageNo, firstPage)
		}

		nav = append(nav, fmt.Sprintf("Prev (%d)", prevPage))
		pageNo = append(pageNo, prevPage)
	}

	if page != lastPage {
		var nextPage int32 = page + 1

		nav = append(nav, fmt.Sprintf("Next (%d)", nextPage))
		pageNo = append(pageNo, nextPage)

		if nextPage != lastPage {
			nav = append(nav, fmt.Sprintf("Last (%d)", lastPage))
			pageNo = append(pageNo, lastPage)
		}
	}

	if pages > 4 {
		nav = append(nav, "Other")
		pageNo = append(pageNo, PAGINATION_OTHER)
	}

	nav = append(nav, "Quit")
	pageNo = append(pageNo, PAGINATION_QUIT)

	return nav, pageNo
}
