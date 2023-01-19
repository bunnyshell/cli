package organization

import (
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	command := &cobra.Command{
		Use: "show",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := lib.GetContext()
			defer cancel()

			request := lib.GetAPI().OrganizationApi.OrganizationView(ctx, settings.Profile.Context.Organization)

			model, resp, err := request.Execute()

			return lib.FormatRequestResult(cmd, model, resp, err)
		},
	}

	flags := command.Flags()

	idFlag := options.Organization.GetFlag("id")
	flags.AddFlag(idFlag)
	_ = command.MarkFlagRequired(idFlag.Name)

	mainCmd.AddCommand(command)
}
