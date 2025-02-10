package action

import (
	"bunnyshell.com/sdk"
)

type DownParameters struct {
	Resource sdk.ComponentResourceItem
}

type Down struct {
	Action
}

func NewDown(
	environment sdk.EnvironmentItem,
) *Down {
	return &Down{
		Action: *NewAction(environment),
	}
}

func (down *Down) Run(parameters *DownParameters) error {
	debugCmp, err := down.Action.GetDebugCmp(parameters.Resource)
	if err != nil {
		return err
	}

	return debugCmp.Down()
}
