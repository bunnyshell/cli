package lib

import (
	"errors"
	"fmt"

	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/interactive"
	"github.com/spf13/cobra"
)

const (
	PaginationQuit  = -1
	PaginationOther = -2

	ShowOtherMinPages = 4
)

var (
	errHandled = error(nil)
	errQuit    = errors.New("quit")
)

type ModelWithPagination interface {
	GetPage() int32

	GetItemsPerPage() int32

	GetTotalItems() int32
}

type Options interface {
	SetPage(int32)
}

type CollectionGenerator func() (ModelWithPagination, error)

func ShowCollection(cmd *cobra.Command, options Options, generator CollectionGenerator) error {
	var page int32

	for {
		model, err := generator()
		if err != nil {
			return err
		}

		if err = FormatCommandData(cmd, model); err != nil {
			return err
		}

		page, err = interactivePagination(cmd, model)
		if err != nil {
			if errors.Is(err, errQuit) {
				return nil
			}

			return err
		}

		options.SetPage(page)
	}
}

func interactivePagination(cmd *cobra.Command, model ModelWithPagination) (int32, error) {
	if !config.GetSettings().IsStylish() {
		return 0, errQuit
	}

	if config.GetSettings().NonInteractive {
		return 0, errQuit
	}

	navPage, err := ProcessPagination(cmd, model)
	if err != nil {
		return 0, err
	}

	if navPage == PaginationQuit {
		return 0, errQuit
	}

	return navPage, nil
}

func ProcessPagination(cmd *cobra.Command, m ModelWithPagination) (int32, error) {
	page := m.GetPage()
	pages := 1 + (m.GetTotalItems()-1)/m.GetItemsPerPage()

	if page > pages {
		return PaginationQuit, nil
	}

	if pages == 1 {
		return PaginationQuit, nil
	}

	nav, pageNo := getPaginationOptions(page, pages)

	index, _, err := interactive.Choose("Navigate to a different page?", nav)
	if err != nil {
		return PaginationQuit, err
	}

	target := pageNo[index]
	if target == PaginationOther {
		target, err = interactive.AskInt32("Go to page:", interactive.AssertBetween(1, pages))
		if err != nil {
			return PaginationQuit, err
		}
	}

	return target, nil
}

func getPaginationOptions(page int32, pages int32) ([]string, []int32) {
	nav := []string{}
	pageNo := []int32{}

	var (
		firstPage int32 = 1
		lastPage        = pages
	)

	if page != firstPage {
		prevPage := page - 1

		if prevPage != firstPage {
			nav = append(nav, fmt.Sprintf("First (%d)", firstPage))
			pageNo = append(pageNo, firstPage)
		}

		nav = append(nav, fmt.Sprintf("Prev (%d)", prevPage))
		pageNo = append(pageNo, prevPage)
	}

	if page != lastPage {
		nextPage := page + 1

		nav = append(nav, fmt.Sprintf("Next (%d)", nextPage))
		pageNo = append(pageNo, nextPage)

		if nextPage != lastPage {
			nav = append(nav, fmt.Sprintf("Last (%d)", lastPage))
			pageNo = append(pageNo, lastPage)
		}
	}

	if pages > ShowOtherMinPages {
		nav = append(nav, "Other")
		pageNo = append(pageNo, PaginationOther)
	}

	nav = append(nav, "Quit")
	pageNo = append(pageNo, PaginationQuit)

	return nav, pageNo
}
