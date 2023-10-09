package common

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var componentsList string

type PartialActionOptions struct {
	ActionOptions

	IsPartial  bool
	Components []string
}

func NewPartialActionOptions(id string, isPartial bool, components []string) *PartialActionOptions {
	return &PartialActionOptions{
		ActionOptions: *NewActionOptions(id),

		IsPartial:  isPartial,
		Components: components,
	}
}

func (pao *PartialActionOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	pao.ActionOptions.UpdateFlagSet(flags)

	flags.StringVar(&componentsList, "components", componentsList, "Execute a partial action with the set components (comma separated names)")
}

func (pao *PartialActionOptions) ProcessCommand(cmd *cobra.Command) {
	pao.IsPartial = cmd.Flags().Lookup("components").Changed
	pao.Components = strings.Split(componentsList, ",")
}
