package progress

import (
	"time"

	"bunnyshell.com/cli/pkg/api/pipeline"
	"bunnyshell.com/cli/pkg/net"
	"bunnyshell.com/sdk"
)

func EventToPipeline(event *sdk.EventItem, options *Options) (*sdk.PipelineItem, error) {
	resume := net.PauseSpinner()
	defer resume()

	spinner := net.MakeSpinner()

	spinner.Start()
	defer spinner.Stop()

	listOptions := pipeline.NewListOptions()
	listOptions.Event = event.GetId()

	for {
		collection, err := pipeline.List(listOptions)
		if err != nil {
			return nil, err
		}

		if !collection.HasEmbedded() {
			time.Sleep(options.Interval)

			continue
		}

		itemOptions := pipeline.NewItemOptions(collection.Embedded.GetItem()[0].GetId())

		return pipeline.Get(itemOptions)
	}
}
