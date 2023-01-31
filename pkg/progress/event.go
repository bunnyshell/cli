package progress

import (
	"errors"
	"fmt"
	"time"

	"bunnyshell.com/cli/pkg/api/pipeline"
	"bunnyshell.com/cli/pkg/net"
	"bunnyshell.com/sdk"
)

var errUnknownEventType = errors.New("unhandleable event type")

func PipelineFromEvent(event sdk.EventItem, options *Options) error {
	if event.GetType() == "env_deplooy" {
		return Event(event.GetId(), options)
	}

	return fmt.Errorf("%w: %s", errUnknownEventType, event.GetType())
}

func Event(eventID string, options *Options) error {
	if options == nil {
		options = NewOptions()
	}

	resume := net.PauseSpinner()
	defer resume()

	item, err := eventToPipline(eventID, options)
	if err != nil {
		return err
	}

	fmt.Printf("Deploy Pipeline: %s\n", item.GetId())

	return progress(*options, generatorFromID(item.GetId()))
}

func eventToPipline(eventID string, options *Options) (*sdk.PipelineItem, error) {
	spinner := net.MakeSpinner()

	spinner.Start()
	defer spinner.Stop()

	listOptions := pipeline.NewListOptions()
	listOptions.Event = eventID

	for {
		collection, err := pipeline.List(listOptions)
		if err != nil {
			return nil, err
		}

		if collection.Embedded == nil {
			time.Sleep(options.Interval)

			continue
		}

		itemOptions := pipeline.NewItemOptions(collection.Embedded.GetItem()[0].GetId())

		return pipeline.Get(itemOptions)
	}
}
