package progress

import (
	"fmt"
	"time"

	pstatus "bunnyshell.com/cli/pkg/api/pipeline/status"
	"bunnyshell.com/sdk"
	"github.com/briandowns/spinner"
)

type PipelineSyncer func() (*sdk.WorkflowItem, error)

type Progress struct {
	Options Options

	spinner *spinner.Spinner
}

type Options struct {
	Interval time.Duration
}

func NewOptions() *Options {
	return &Options{
		Interval: defaultSpinnerUpdate,
	}
}

func NewPipeline(options Options) *Progress {
	spinner := spinner.New(spinner.CharSets[defaultProgressSet], defaultSpinnerUpdate)
	spinner.Prefix = fmt.Sprintf(
		"%s Fetching pipeline status... ",
		statusMap[PipelineWorking],
	)

	return &Progress{
		Options: options,

		spinner: spinner,
	}
}

func (p *Progress) Update(pipelineSync PipelineSyncer) error {
	for {
		workflow, err := pipelineSync()
		if err != nil {
			return err
		}

		waiting, err := p.UpdatePipeline(workflow)
		if err != nil {
			return err
		}

		if !waiting {
			return nil
		}

		time.Sleep(p.Options.Interval)
	}
}

func (p *Progress) UpdatePipeline(workflow *sdk.WorkflowItem) (bool, error) {
	if workflow == nil {
		return false, nil
	}

	p.spinner.Prefix = fmt.Sprintf(
		"%s Processing... %d/%d jobs completed ",
		statusMap[p.getState(workflow.GetStatus())],
		workflow.GetCompletedJobsCount(),
		workflow.GetJobsCount(),
	)

	switch workflow.GetStatus() {
	case pstatus.WorkflowInProgress, pstatus.WorkflowAborting, pstatus.WorkflowFailing, pstatus.WorkflowThrottled, pstatus.WorkflowQueued:
		return true, nil
	case pstatus.WorkflowSuccess:
		return false, nil
	case pstatus.WorkflowFailed:
		return false, ErrPipeline
	case pstatus.WorkflowAborted:
		return false, ErrPipelineAborted
	default:
		return false, fmt.Errorf("%w: unknown status %s", ErrPipeline, workflow.GetStatus())
	}
}

func (p *Progress) Start() {
	p.spinner.Start()
}

func (p *Progress) Stop() {
	p.spinner.Stop()
}

func (p *Progress) getState(status string) PipelineStatus {
	switch status {
	case pstatus.WorkflowSuccess:
		return PipelineFinished
	case pstatus.WorkflowInProgress, pstatus.WorkflowAborting, pstatus.WorkflowFailing, pstatus.WorkflowThrottled, pstatus.WorkflowQueued:
		return PipelineWorking
	case pstatus.WorkflowFailed:
		return PipelineFailed
	case pstatus.WorkflowAborted:
		return PipelineAborted
	}

	return PipelineUnknownState
}
