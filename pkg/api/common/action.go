package common

import (
	"time"

	"github.com/spf13/pflag"
)

type ActionOptions struct {
	ItemOptions

	WithoutPipeline bool

	Interval time.Duration
}

func NewActionOptions(id string) *ActionOptions {
	return &ActionOptions{
		ItemOptions: *NewItemOptions(id),

		WithoutPipeline: false,

		Interval: 2000 * time.Millisecond,
	}
}

func (ao *ActionOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	flags.BoolVar(&ao.WithoutPipeline, "no-wait", ao.WithoutPipeline, "Do not wait for pipeline until finish")
	flags.DurationVar(&ao.Interval, "pipeline-monitor-interval", ao.Interval, "Pipeline check interval")
}
