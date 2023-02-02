package common

import (
	"github.com/spf13/pflag"
)

type ListOptions struct {
	Options

	Page int32
}

func NewListOptions() *ListOptions {
	return &ListOptions{
		Page: 1,
	}
}

func (lo *ListOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	flags.Int32Var(&lo.Page, "page", lo.Page, "Listing Page")
}
