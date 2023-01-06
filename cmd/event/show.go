package event

import (
	"net/http"
	"time"

	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
	"github.com/spf13/cobra"
)

func init() {
	var (
		monitor   bool
		eventID   string
		lastEvent *sdk.EventItem
	)

	idleNotify := 10 * time.Second
	errWait := idleNotify / 2

	command := &cobra.Command{
		Use: "show",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			resp, r, err := getEvent(eventID)
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
				resp, _, err := getEvent(eventID)
				if err != nil {
					if now.After(idleThreshold) {
						_ = lib.FormatCommandError(cmd, err)
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
				_ = lib.FormatCommandData(cmd, resp)
				idleThreshold = now.Add(idleNotify)
			}
		},
	}

	command.Flags().StringVar(&eventID, "id", eventID, "Event Id")
	command.MarkFlagRequired("id")

	command.Flags().BoolVar(&monitor, "monitor", false, "monitor the event for changes or until finished")
	command.Flags().DurationVar(&idleNotify, "idle-notify", idleNotify, "Network timeout on requests")

	mainCmd.AddCommand(command)
}

func isFinalStatus(e *sdk.EventItem) bool {
	return e.GetStatus() == "success" || e.GetStatus() == "error"
}

func getEvent(eventID string) (*sdk.EventItem, *http.Response, error) {
	ctx, cancel := lib.GetContext()
	defer cancel()

	request := lib.GetAPI().EventApi.EventView(ctx, eventID)

	return request.Execute()
}
