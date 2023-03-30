package event

import (
	"fmt"
	"time"

	"bunnyshell.com/cli/pkg/api/event"
	"bunnyshell.com/cli/pkg/build"
	"bunnyshell.com/cli/pkg/config/option"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/net"
	"bunnyshell.com/sdk"
	"github.com/spf13/cobra"
)

func init() {
	var monitor bool

	idleNotify := 10 * time.Second

	itemOptions := event.NewItemOptions("")

	command := &cobra.Command{
		Use: "show",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			model, err := event.Get(itemOptions)
			if err != nil {
				return err
			}

			_ = lib.FormatCommandData(cmd, model)

			if !monitor || isFinalStatus(model) {
				return nil
			}

			resume := net.PauseSpinner()
			defer resume()

			spinner := net.MakeSpinner()

			spinner.Start()
			defer spinner.Stop()

			monitorEvent(cmd, model, idleNotify)

			return nil
		},
	}

	flags := command.Flags()

	flags.AddFlag(getIDOption(&itemOptions.ID).GetRequiredFlag("id"))

	flags.BoolVar(&monitor, "monitor", false, "monitor the event for changes or until finished")
	flags.DurationVar(&idleNotify, "idle-notify", idleNotify, "Network timeout on requests")

	mainCmd.AddCommand(command)
}

func getIDOption(value *string) *option.String {
	help := fmt.Sprintf(
		`Find available Events with "%s events list"`,
		build.Name,
	)

	idOption := option.NewStringOption(value)

	idOption.AddFlagWithExtraHelp("id", "Environment Variable Id", help)

	return idOption
}

func isFinalStatus(e *sdk.EventItem) bool {
	switch e.GetStatus() {
	case "success", "error", "delegated":
		return true
	default:
		return false
	}
}

func monitorEvent(cmd *cobra.Command, lastEvent *sdk.EventItem, idleNotify time.Duration) {
	itemOptions := event.NewItemOptions(lastEvent.GetId())

	errWait := idleNotify / 2

	idleThreshold := time.Now().Add(idleNotify)

	for {
		now := time.Now()

		model, err := event.Get(itemOptions)
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

		if lastEvent.GetUpdatedAt().Equal(model.GetUpdatedAt()) {
			continue
		}

		if isFinalStatus(model) {
			return
		}

		lastEvent = model
		_ = lib.FormatCommandData(cmd, model)
		idleThreshold = now.Add(idleNotify)
	}
}
