package pipeline

import (
	"bunnyshell.com/cli/pkg/api/workflow_job"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	listOptions := workflow_job.NewListOptions()

	var pipelineID string
	var jobStatuses []string

	command := &cobra.Command{
		Use: "jobs",

		Short: "List jobs in a pipeline",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			listOptions.Workflow = pipelineID
			listOptions.Status = jobStatuses

			return lib.ShowCollection(cmd, listOptions, func() (lib.ModelWithPagination, error) {
				return workflow_job.List(listOptions)
			})
		},
	}

	flags := command.Flags()

	flags.AddFlag(getIDOption(&pipelineID).GetRequiredFlag("id"))
	flags.StringArrayVar(&jobStatuses, "jobStatus", jobStatuses, "Filter by Job Status (repeatable)")

	listOptions.UpdateFlagSet(flags)

	mainCmd.AddCommand(command)
}
