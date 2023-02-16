package progress

import (
	"fmt"
	"time"

	"bunnyshell.com/cli/pkg/api/event"
	"bunnyshell.com/sdk"
	"github.com/briandowns/spinner"
)

func handleWorkflow(eventItem *sdk.EventItem, options *Options, spinner *spinner.Spinner) (*sdk.EventItem, error) {
	for {
		if eventItem.GetType() != "env_queued" {
			return eventItem, nil
		}

		delegatedEvent, err := delegateEvent(eventItem.GetId(), options)
		if err != nil {
			return nil, err
		}

		spinner.FinalMSG = fmt.Sprintf("Delegating to EventID: %s\n", delegatedEvent.GetId())
		spinner.Restart()
		spinner.FinalMSG = ""

		eventItem = delegatedEvent
	}
}

func delegateEvent(eventID string, options *Options) (*sdk.EventItem, error) {
	itemOptions := event.NewItemOptions(eventID)

	for {
		eventItem, err := event.Get(itemOptions)
		if err != nil {
			return nil, err
		}

		switch eventItem.GetStatus() {
		case "fail":
			return nil, fmt.Errorf("event has failed")

		case "aborted":
			return nil, fmt.Errorf("event was aborted")

		case "delegated":
			delegatedID := eventItem.GetDelegated()

			return event.Get(event.NewItemOptions(delegatedID))

		default:
			time.Sleep(options.Interval)
		}
	}
}
