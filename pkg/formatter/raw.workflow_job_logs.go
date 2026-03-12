package formatter

import (
	"bytes"
	"fmt"
)

func rawWorkflowJobLogs(data []WorkflowJobLogsResult) []byte {
	var buf bytes.Buffer

	for i, result := range data {
		if result.Logs == nil {
			continue
		}

		if i > 0 {
			fmt.Fprintf(&buf, "\n\n")
		}

		// Job header
		jobName := result.JobID
		if job, ok := result.Logs.GetWorkflowJobOk(); ok && job != nil {
			jobName = getJobDisplayName(job, result.JobID)
		}
		fmt.Fprintf(&buf, "Job: %s\n", jobName)

		// Steps
		if steps, ok := result.Logs.GetStepsOk(); ok && steps != nil {
			for j, step := range steps {
				if j > 0 {
					fmt.Fprintf(&buf, "\n")
				}

				stepName := "Unknown"
				if name, ok := step.GetNameOk(); ok && name != nil {
					stepName = *name
				}
				fmt.Fprintf(&buf, "Step: %s\n", stepName)

				if logs, ok := step.GetLogsOk(); ok && logs != nil {
					for _, logEntry := range logs {
						timestamp := ""
						if t, ok := logEntry.GetTimeOk(); ok && t != nil {
							timestamp = t.Format("15:04:05.000")
						}

						message := ""
						if msg, ok := logEntry.GetLogOk(); ok && msg != nil {
							message = *msg
						}

						if timestamp != "" {
							fmt.Fprintf(&buf, "%s  %s\n", timestamp, message)
						} else {
							fmt.Fprintf(&buf, "%s\n", message)
						}
					}
				}
			}
		}
	}

	return buf.Bytes()
}
