package organization

import (
	"bunnyshell.com/cli/pkg/api/organization"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	listOptions := organization.NewListOptions()

	command := &cobra.Command{
		Use: "list",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			return lib.ShowCollection(cmd, listOptions, func() (lib.ModelWithPagination, error) {
				return organization.List(listOptions)
			})
		},
	}

	flags := command.Flags()

	listOptions.UpdateFlagSet(flags)

	mainCmd.AddCommand(command)
}
