package workflow_job

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
)

// PipelineLogs wraps logs from all jobs in a workflow
type PipelineLogs struct {
	WorkflowID string            `json:"workflowId"`
	Jobs       []WorkflowJobLogs `json:"jobs"`
}

// JobInfo contains metadata about a workflow job
type JobInfo struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Status string `json:"status"`
}

// WorkflowJobLogs represents the structure of workflow job logs API response
type WorkflowJobLogs struct {
	WorkflowJobID string     `json:"workflowJobId"`
	JobName       string     `json:"jobName,omitempty"`
	Status        string     `json:"status"`
	Steps         []LogStep  `json:"steps"`
	Pagination    Pagination `json:"pagination"`
}

// LogStep represents a single step in the workflow
type LogStep struct {
	Name      string       `json:"name"`
	Status    string       `json:"status"`
	ExitCode  int          `json:"exitCode"`
	StartedAt string       `json:"startedAt"`
	Logs      []LogMessage `json:"logs"`
}

// LogMessage represents a single log message
type LogMessage struct {
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}

// Pagination contains pagination metadata
type Pagination struct {
	Offset  int  `json:"offset"`
	Limit   int  `json:"limit"`
	Total   int  `json:"total"`
	HasMore bool `json:"hasMore"`
}

// LogsOptions contains options for fetching workflow job logs
type LogsOptions struct {
	Profile config.Profile
	JobID   string
	Offset  int
	Limit   int
}

// GetLogs fetches workflow job logs from the API
func GetLogs(options *LogsOptions) (*WorkflowJobLogs, error) {
	ctx, cancel := lib.GetContextFromProfile(options.Profile)
	defer cancel()

	// Build API URL
	apiURL := buildAPIURL(options.Profile, options.JobID, options.Offset, options.Limit)

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authorization header
	req.Header.Set("X-Auth-Token", options.Profile.Token)
	req.Header.Set("Accept", "application/hal+json")

	if config.GetSettings().Debug {
		fmt.Fprintf(os.Stderr, "GET %s\n", apiURL)
	}

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return nil, parseHTTPError(resp.StatusCode, body)
	}

	// Parse JSON response
	var logs WorkflowJobLogs
	if err := json.Unmarshal(body, &logs); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &logs, nil
}

// buildAPIURL constructs the full API URL with query parameters
func buildAPIURL(profile config.Profile, jobID string, offset, limit int) string {
	scheme := profile.Scheme
	if scheme == "" {
		scheme = "https"
	}
	baseURL := fmt.Sprintf("%s://%s", scheme, profile.Host)
	path := fmt.Sprintf("/v1/workflow_jobs/%s/logs", jobID)

	// Build query parameters
	params := url.Values{}
	params.Add("offset", fmt.Sprintf("%d", offset))
	params.Add("limit", fmt.Sprintf("%d", limit))

	return fmt.Sprintf("%s%s?%s", baseURL, path, params.Encode())
}

// parseHTTPError creates a user-friendly error message from HTTP response
func parseHTTPError(statusCode int, body []byte) error {
	var errorResp struct {
		Error string `json:"error"`
		Code  string `json:"code"`
	}

	// Try to parse error response
	if err := json.Unmarshal(body, &errorResp); err == nil && errorResp.Error != "" {
		return fmt.Errorf("%s (HTTP %d)", errorResp.Error, statusCode)
	}

	// Fallback to generic error messages
	switch statusCode {
	case http.StatusNotFound:
		return fmt.Errorf("workflow job not found (HTTP 404)")
	case http.StatusUnauthorized:
		return fmt.Errorf("authentication failed. Run 'bns configure' to set your token (HTTP 401)")
	case http.StatusForbidden:
		return fmt.Errorf("access forbidden. You don't have permission to view these logs (HTTP 403)")
	case http.StatusTooManyRequests:
		return fmt.Errorf("rate limit exceeded. Please wait and try again (HTTP 429)")
	case http.StatusBadGateway:
		return fmt.Errorf("unable to retrieve log file from storage. Please try again later (HTTP 502)")
	default:
		return fmt.Errorf("API error (HTTP %d): %s", statusCode, string(body))
	}
}

// FetchAllPages fetches all pages of logs automatically
func FetchAllPages(options *LogsOptions) (*WorkflowJobLogs, error) {
	var allLogs *WorkflowJobLogs
	var allSteps []LogStep

	offset := options.Offset
	limit := options.Limit

	for {
		// Fetch current page
		opts := &LogsOptions{
			Profile: options.Profile,
			JobID:   options.JobID,
			Offset:  offset,
			Limit:   limit,
		}

		logs, err := GetLogs(opts)
		if err != nil {
			return nil, err
		}

		// Store first page metadata
		if allLogs == nil {
			allLogs = logs
			allSteps = logs.Steps
		} else {
			// Merge steps from subsequent pages
			allSteps = mergeSteps(allSteps, logs.Steps)
		}

		// Check if more pages exist
		if !logs.Pagination.HasMore {
			break
		}

		// Move to next page
		offset += limit
	}

	// Update final result
	if allLogs != nil {
		allLogs.Steps = allSteps
		allLogs.Pagination.HasMore = false
	}

	return allLogs, nil
}

// GetJobInfo fetches metadata (name, type, status) for a workflow job
func GetJobInfo(profile config.Profile, jobID string) (*JobInfo, error) {
	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	scheme := profile.Scheme
	if scheme == "" {
		scheme = "https"
	}
	apiURL := fmt.Sprintf("%s://%s/v1/workflow_jobs/%s", scheme, profile.Host, jobID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Auth-Token", profile.Token)
	req.Header.Set("Accept", "application/hal+json")

	if config.GetSettings().Debug {
		fmt.Fprintf(os.Stderr, "GET %s\n", apiURL)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch job info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read job info response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, parseHTTPError(resp.StatusCode, body)
	}

	var info JobInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("failed to parse job info: %w", err)
	}

	return &info, nil
}

// mergeSteps merges log steps, combining logs from the same step
func mergeSteps(existing, new []LogStep) []LogStep {
	// If last existing step matches first new step, merge their logs
	if len(existing) > 0 && len(new) > 0 {
		lastExisting := &existing[len(existing)-1]
		firstNew := new[0]

		if lastExisting.Name == firstNew.Name {
			// Merge logs
			lastExisting.Logs = append(lastExisting.Logs, firstNew.Logs...)

			// Append remaining new steps
			return append(existing, new[1:]...)
		}
	}

	// No overlap, just append
	return append(existing, new...)
}
