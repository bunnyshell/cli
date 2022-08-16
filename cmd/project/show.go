package project

import (
	"github.com/spf13/cobra"

	"bunnyshell.com/cli/pkg/lib"
)

func init() {
	project := &lib.CLIContext.Profile.Context.Project

	command := &cobra.Command{
		Use: "show",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := lib.GetContext()
			defer cancel()

			request := lib.GetAPI().ProjectApi.ProjectView(ctx, *project)

			resp, r, err := request.Execute()

			return lib.FormatRequestResult(cmd, resp, r, err)
		},
	}

	command.Flags().StringVar(project, "id", *project, "Project Id")
	command.MarkFlagRequired("id")

	mainCmd.AddCommand(command)
}
