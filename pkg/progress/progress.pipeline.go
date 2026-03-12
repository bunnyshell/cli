package progress

import (
	"fmt"
	"time"

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
	case StatusInProgress, StatusPending:
		return true, nil
	case StatusSuccess:
		return false, nil
	case StatusFailed:
		return false, ErrPipeline
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
	case StatusSuccess:
		return PipelineFinished
	case StatusInProgress, StatusPending:
		return PipelineWorking
	case StatusFailed:
		return PipelineFailed
	}

	return PipelineUnknownState
}
