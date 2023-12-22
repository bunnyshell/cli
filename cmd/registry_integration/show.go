package registry_integration

import (
	"fmt"

	"bunnyshell.com/cli/pkg/api/registry_integration"
	"bunnyshell.com/cli/pkg/build"
	"bunnyshell.com/cli/pkg/config/option"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	itemOptions := registry_integration.NewItemOptions("")

	command := &cobra.Command{
		Use:     "show",
		GroupID: mainGroup.ID,

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			model, err := registry_integration.Get(itemOptions)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			return lib.FormatCommandData(cmd, model)
		},
	}

	flags := command.Flags()

	flags.AddFlag(getIDOption(&itemOptions.ID).GetRequiredFlag("id"))

	mainCmd.AddCommand(command)
}

func getIDOption(value *string) *option.String {
	help := fmt.Sprintf(
		`Find available Container Registries Integrations with "%s container-registries list"`,
		build.Name,
	)

	idOption := option.NewStringOption(value)

	idOption.AddFlagWithExtraHelp("id", "Container Registry Integration Id", help)

	return idOption
}
