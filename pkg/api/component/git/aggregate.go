package git

import (
	"bunnyshell.com/sdk"
	"github.com/spf13/pflag"
)

type AggregateOptions struct {
	ListOptions

	Limit int
}

func NewAggregateOptions() *AggregateOptions {
	return &AggregateOptions{
		ListOptions: *NewListOptions(),

		Limit: 0,
	}
}

func (ao *AggregateOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	flags.IntVar(&ao.Limit, "limit", ao.Limit, "Limit aggregated results, 0=no limit")

	// There's no pagination when aggregating
	ao.ListOptions.updateSelfFlags(flags)
}

func Aggregate(options *AggregateOptions) ([]sdk.ComponentGitCollection, error) {
	result := []sdk.ComponentGitCollection{}

	for {
		model, err := List(&options.ListOptions)
		if err != nil {
			return nil, err
		}

		result = append(result, getPaginatedGitComponents(model)...)

		if !hasNextPage(model) {
			return result, nil
		}

		if options.Limit != 0 && len(result) > options.Limit {
			return result[0:options.Limit], nil
		}

		options.Page++
	}
}

func getPaginatedGitComponents(model *sdk.PaginatedComponentGitCollection) []sdk.ComponentGitCollection {
	if !model.HasEmbedded() {
		return []sdk.ComponentGitCollection{}
	}

	return model.Embedded.Item
}

func hasNextPage(model *sdk.PaginatedComponentGitCollection) bool {
	return model.HasLinks() && model.Links.HasNext()
}
