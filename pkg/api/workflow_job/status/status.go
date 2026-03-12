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

// Job status values.
const (
	JobPending     = "pending"
	JobQueued      = "queued"
	JobInProgress  = "in_progress"

	JobFailed      = "failed"
	JobAbortFailed = "abort_failed"
	JobSuccess     = "success"

	JobSkipped  = "skipped"
	JobAborting = "aborting"
	JobAborted  = "aborted"
)

// Step status values.
const (
	StepFailed  = "failed"
	StepSuccess = "success"
)
