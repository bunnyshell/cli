package status

// Workflow status values.
const (
	WorkflowQueued     = "queued"
	WorkflowThrottled  = "throttled"
	WorkflowInProgress = "in_progress"
	WorkflowSuccess    = "success"

	WorkflowFailing = "failing"
	WorkflowFailed  = "failed"

	WorkflowAborting = "aborting"
	WorkflowAborted  = "aborted"
)
