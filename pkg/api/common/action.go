package common

import (
	"github.com/spf13/pflag"
)

type ActionOptions struct {
	ItemOptions

	WithoutPipeline bool
}

func NewActionOptions(id string) *ActionOptions {
	return &ActionOptions{
		ItemOptions: *NewItemOptions(id),

		WithoutPipeline: false,
	}
}

func (ao *ActionOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	flags.BoolVar(&ao.WithoutPipeline, "no-wait", ao.WithoutPipeline, "Do not wait for pipeline until finish")
}
