package action

import (
	"bunnyshell.com/sdk"
)

type DownParameters struct {
	Resource sdk.ComponentResourceItem

	OverrideClusterServer string
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
	remoteDev, err := down.Action.GetRemoteDev(parameters.Resource, parameters.OverrideClusterServer)
	if err != nil {
		return err
	}

	return remoteDev.Down()
}
