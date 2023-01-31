package progress

import (
	"bunnyshell.com/cli/pkg/api/pipeline"
	"bunnyshell.com/cli/pkg/net"
	"bunnyshell.com/sdk"
)

func Pipeline(pipelineID string, options *Options) error {
	if options == nil {
		options = NewOptions()
	}

	resume := net.PauseSpinner()
	defer resume()

	return progress(*options, generatorFromID(pipelineID))
}

func progress(options Options, generate PipelineSyncer) error {
	wizard := NewPipeline(options)

	wizard.Start()
	defer wizard.Stop()

	return wizard.Update(generate)
}

func generatorFromID(id string) PipelineSyncer {
	itemOptions := pipeline.NewItemOptions(id)

	return func() (*sdk.PipelineItem, error) {
		return pipeline.Get(itemOptions)
	}
}
