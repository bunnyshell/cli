package common

import (
	"github.com/spf13/pflag"
)

type ActionOptions struct {
	ItemOptions

	WithPipeline bool
}

func NewActionOptions(id string) *ActionOptions {
	return &ActionOptions{
		ItemOptions: *NewItemOptions(id),

		WithPipeline: true,
	}
}

func (ao *ActionOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	flags.BoolVar(&ao.WithPipeline, "no-wait", ao.WithPipeline, "Do not wait for pipeline until finish")
}
