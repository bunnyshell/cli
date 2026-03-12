package formatter

import "errors"

var errRawUnsupported = errors.New("raw output is not supported for this command")

func raw(data interface{}) ([]byte, error) {
	switch dataType := data.(type) {
	case []WorkflowJobLogsResult:
		return rawWorkflowJobLogs(dataType), nil
	}

	return nil, errRawUnsupported
}
