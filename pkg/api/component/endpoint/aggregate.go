package endpoint

import (
	"bunnyshell.com/cli/pkg/net"
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

func Aggregate(options *AggregateOptions) ([]sdk.ComponentEndpointCollection, error) {
	result := []sdk.ComponentEndpointCollection{}

	resume := net.PauseSpinner()
	defer resume()

	spinner := net.MakeSpinner()

	spinner.Start()
	defer spinner.Stop()

	for {
		model, err := List(&options.ListOptions)
		if err != nil {
			return nil, err
		}

		result = append(result, getPaginatedComponents(model)...)

		if !hasNextPage(model) {
			return result, nil
		}

		if options.Limit != 0 && len(result) > options.Limit {
			return result[0:options.Limit], nil
		}

		options.Page++
	}
}

func getPaginatedComponents(model *sdk.PaginatedComponentEndpointCollection) []sdk.ComponentEndpointCollection {
	if !model.HasEmbedded() {
		return []sdk.ComponentEndpointCollection{}
	}

	return model.Embedded.Item
}

func hasNextPage(model *sdk.PaginatedComponentEndpointCollection) bool {
	return model.HasLinks() && model.Links.HasNext()
}
