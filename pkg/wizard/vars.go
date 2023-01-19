package wizard

import (
	"errors"
)

const (
	selectPage = iota

	paginationItemsPerPage = 30
)

var (
	paginationOptions = map[int]string{
		selectPage: "Select Page",
	}

	paginationDisplaySize = paginationItemsPerPage + len(paginationOptions)

	ErrEmptyListing = errors.New("no resources found")
)
