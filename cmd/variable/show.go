package variable

import (
	"github.com/spf13/cobra"

	"bunnyshell.com/cli/pkg/lib"
)

func init() {
	var id string

	command := &cobra.Command{
		Use: "show",
		RunE: func(cmd *cobra.Command, args []string) error {
			var api = lib.GetAPI().EnvironmentVariableApi

			ctx, cancel := lib.GetContext()
			defer cancel()

			request := api.EnvironmentVariableView(ctx, id)

			resp, r, err := request.Execute()
			return lib.FormatRequestResult(cmd, resp, r, err)
		},
	}

	command.Flags().StringVar(&id, "id", id, "Environment Variable Id")
	command.MarkFlagRequired("id")

	mainCmd.AddCommand(command)
}
