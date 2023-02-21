package progress

import (
	"fmt"
	"time"

	"bunnyshell.com/sdk"
	"github.com/briandowns/spinner"
)

type PipelineSyncer func() (*sdk.PipelineItem, error)

type Progress struct {
	Options Options

	spinner *spinner.Spinner

	stages map[string]bool
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
		stages:  map[string]bool{},
	}
}

func (p *Progress) Update(pipelineSync PipelineSyncer) error {
	for {
		pipeline, err := pipelineSync()
		if err != nil {
			return err
		}

		waiting, err := p.UpdatePipeline(pipeline)
		if err != nil {
			return err
		}

		if !waiting {
			return nil
		}

		time.Sleep(p.Options.Interval)
	}
}

func (p *Progress) UpdatePipeline(pipeline *sdk.PipelineItem) (bool, error) {
	if pipeline == nil {
		return false, nil
	}

	p.spinner.Prefix = "Processing Pipeline "

	for _, stage := range pipeline.GetStages() {
		switch p.setStage(stage) {
		case Success:
			continue
		case Failed:
			return false, ErrPipeline
		case Synced:
			return true, nil
		}
	}

	switch pipeline.GetStatus() {
	case StatusInProgress, StatusPending:
		return true, nil
	case StatusSuccess:
		return false, nil
	case StatusFailed:
		return false, ErrPipeline
	default:
		return false, fmt.Errorf("%w: unknown status %s", ErrPipeline, pipeline.GetStatus())
	}
}

func (p *Progress) Start() {
	p.spinner.Start()
}

func (p *Progress) Stop() {
	p.spinner.Stop()
}

func (p *Progress) setStage(stage sdk.StageItem) UpdateStatus {
	if stage.GetStatus() == StatusFailed {
		p.finishStage(stage)

		return Failed
	}

	if stage.GetStatus() == StatusSuccess {
		p.finishStage(stage)

		return Success
	}

	p.syncStage(stage)

	return Synced
}

func (p *Progress) finishStage(stage sdk.StageItem) {
	if p.stages[stage.GetId()] {
		return
	}

	p.stages[stage.GetId()] = true

	p.spinner.FinalMSG = fmt.Sprintf(
		"%s %s finished %d jobs in %s\n",
		statusMap[p.getState(stage)],
		stage.GetName(),
		stage.GetJobsCount(),
		time.Duration(stage.GetDuration())*time.Second,
	)

	p.spinner.Restart()

	p.spinner.FinalMSG = ""
}

func (p *Progress) syncStage(stage sdk.StageItem) {
	p.spinner.Prefix = fmt.Sprintf(
		"%s %s... %d/%d jobs completed ",
		statusMap[p.getState(stage)],
		stage.GetName(),
		stage.GetCompletedJobsCount(),
		stage.GetJobsCount(),
	)
}

func (p *Progress) getState(stage sdk.StageItem) PipelineStatus {
	switch stage.GetStatus() {
	case StatusSuccess:
		return PipelineFinished
	case StatusInProgress, StatusPending:
		return PipelineWorking
	case StatusFailed:
		return PipelineFailed
	}

	return PipelineUnknownState
}
