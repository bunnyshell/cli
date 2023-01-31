package k8sIntegration

import (
	"bunnyshell.com/cli/pkg/api/k8s"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	itemOptions := k8s.NewItemOptions("")

	command := &cobra.Command{
		Use:     "show",
		GroupID: mainGroup.ID,

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			model, err := k8s.Get(itemOptions)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			return lib.FormatCommandData(cmd, model)
		},
	}

	flags := command.Flags()

	idFlagName := "id"
	flags.StringVar(&itemOptions.ID, idFlagName, itemOptions.ID, "Kubernetes Integrations Id")
	_ = command.MarkFlagRequired(idFlagName)

	mainCmd.AddCommand(command)
}
