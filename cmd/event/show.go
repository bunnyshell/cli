package event

import (
	"net/http"
	"time"

	"github.com/spf13/cobra"

	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
)

func init() {
	var monitor bool
	var eventId string
	var lastEvent *sdk.EventItem
	idleNotify := 10 * time.Second
	errWait := idleNotify / 2

	command := &cobra.Command{
		Use: "show",
		RunE: func(cmd *cobra.Command, args []string) error {
			resp, r, err := getEvent(eventId)
			lastEvent = resp
			return lib.FormatRequestResult(cmd, resp, r, err)
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			if !monitor || isFinalStatus(lastEvent) {
				return
			}

			idleThreshold := time.Now().Add(idleNotify)
			for {
				now := time.Now()
				resp, _, err := getEvent(eventId)

				if err != nil {
					if now.After(idleThreshold) {
						lib.FormatCommandError(cmd, err)
						time.Sleep(errWait)
						idleThreshold = now.Add(idleNotify)
					} else {
						time.Sleep(errWait)
					}
					continue
				}

				if lastEvent.GetUpdatedAt().Equal(resp.GetUpdatedAt()) {
					continue
				}

				if isFinalStatus(resp) {
					return
				}

				lastEvent = resp
				lib.FormatCommandData(cmd, resp)
				idleThreshold = now.Add(idleNotify)
			}
		},
	}

	command.Flags().StringVar(&eventId, "id", eventId, "Event Id")
	command.MarkFlagRequired("id")

	command.Flags().BoolVar(&monitor, "monitor", false, "monitor the event for changes or until finished")
	command.Flags().DurationVar(&idleNotify, "idle-notify", idleNotify, "Network timeout on requests")
	command.Flags().DurationVar(&idleNotify, "test", idleNotify, "Network timeout on requests")

	mainCmd.AddCommand(command)
}

func isFinalStatus(e *sdk.EventItem) bool {
	return e.GetStatus() == "success" || e.GetStatus() == "error"
}

func getEvent(eventId string) (*sdk.EventItem, *http.Response, error) {
	ctx, cancel := lib.GetContext()
	defer cancel()

	request := lib.GetAPI().EventApi.EventView(ctx, eventId)
	return request.Execute()
}
