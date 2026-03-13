package action

import (
	"net/http"

	"bunnyshell.com/cli/pkg/api/environment"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	abortOptions := environment.NewAbortOptions("")

	command := &cobra.Command{
		Use: "abort",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			abortOptions.ID = settings.Profile.Context.Environment

			event, err := environment.Abort(abortOptions)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			if event != nil {
				return lib.FormatCommandData(cmd, event)
			}

			if settings.IsStylish() {
				cmd.Println("Nothing to abort")
				return nil
			}

			return lib.FormatCommandData(cmd, map[string]interface{}{
				"status": http.StatusNoContent,
				"detail": "Nothing to abort",
			})
		},
	}

	flags := command.Flags()
	flags.AddFlag(options.Environment.GetRequiredFlag("id"))

	mainCmd.AddCommand(command)
}
