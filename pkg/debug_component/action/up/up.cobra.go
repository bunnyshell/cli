package up

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func (up *Options) UpdateFlagSet(
	command *cobra.Command,
	flags *pflag.FlagSet,
) {
	flags.BoolVar(
		&up.ManualSelectSingleResource,
		"manual-select-single-resource",
		up.ManualSelectSingleResource,
		"Do not skip interactive selectors when there is only one result",
	)

	_ = flags.MarkHidden("no-auto-select-one")

	flags.BoolVar(
        &up.ForceRecreateResource,
        "force-recreate-resource",
        up.ForceRecreateResource,
        "Force recreate Pod even if another debug session is in progress. May break the debug down command",
    )

	flags.DurationVarP(&up.waitTimeout, "wait-timeout", "w", up.waitTimeout, "Time to wait for the pod to be ready")

	flags.StringVarP(&up.resourcePath, "resource", "s", up.resourcePath, "The cluster resource to use (namespace/kind/name format).")
	flags.StringVar(&up.containerName, "container", up.containerName, "The container name to use for remote development")

	up.addContainerConfigFlags(flags)
}

func (up *Options) addContainerConfigFlags(
	flags *pflag.FlagSet,
) {
	flags.StringArrayVar(&up.environPairs, "env", up.environPairs, "Environment variables to set in the debugged container")

	flags.StringVar(&up.limitCPU, "limit-cpu", up.limitCPU, "CPU resource limits for the debugged container")
	flags.StringVar(&up.limitMemory, "limit-memory", up.limitMemory, "Memory resource limits for the debugged container")
	flags.StringVar(&up.requestCPU, "request-cpu", up.requestCPU, "CPU resource reservations for the debugged container")
	flags.StringVar(&up.requestMemory, "request-memory", up.requestMemory, "Memory resource reservations for the debugged container")
}
