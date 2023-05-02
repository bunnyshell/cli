package down

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func (down *Options) UpdateFlagSet(
	command *cobra.Command,
	flags *pflag.FlagSet,
) {
	flags.StringVarP(&down.resourcePath, "resource", "s", down.resourcePath, "The cluster resource to use (namespace/kind/name format).")

	down.manager.UpdateFlagSet(command, flags)
}
