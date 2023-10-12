package common

import (
	"fmt"

	"github.com/spf13/pflag"
)

const componentVarName = "component"

type PartialActionOptions struct {
	ActionOptions

	components []string

	flags *pflag.FlagSet
}

func NewPartialActionOptions(id string) *PartialActionOptions {
	return &PartialActionOptions{
		ActionOptions: *NewActionOptions(id),

		components: []string{},
	}
}

func (pao *PartialActionOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	pao.ActionOptions.UpdateFlagSet(flags)

	// We'll need this later to check if the flag was provided at all - in order to determine if it was a partial action or not.
	pao.flags = flags

	usage := fmt.Sprintf("Execute a partial action with the set components. Provide \"--%s ''\" for a no-components operation.", componentVarName)
	flags.StringArrayVar(&pao.components, componentVarName, pao.components, usage)
}

// Handle the "--component ”" case.
func (pao *PartialActionOptions) GetActionComponents() []string {
	if len(pao.components) == 1 && pao.components[0] == "" {
		return []string{}
	}

	return pao.components
}

func (pao *PartialActionOptions) IsPartial() bool {
	if pao.flags == nil {
		return false
	}

	return pao.flags.Lookup(componentVarName).Changed
}
