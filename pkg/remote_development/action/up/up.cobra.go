package up

import (
	"fmt"
	"strings"

	remoteDevMutagenConfig "bunnyshell.com/dev/pkg/mutagen/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/thediveo/enumflag/v2"
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

	flags.DurationVarP(&up.waitTimeout, "wait-timeout", "w", up.waitTimeout, "Time to wait for the pod to be ready")

	flags.StringVarP(
		&up.localSyncPath,
		"local-sync-path",
		"l",
		up.localSyncPath,
		"The folder on your machine that will be synced into the container on the path specified by --remote-sync-path",
	)
	flags.StringVarP(
		&up.remoteSyncPath,
		"remote-sync-path",
		"r",
		up.remoteSyncPath,
		"The folder within the container where the source code of the application resides\n"+
			"This will be used as a persistent volume to preserve your changes across multiple development sessions\n"+
			"When using --sync-mode=none it will be used only as a workspace where changes to those files will be preserved",
	)

	flags.StringVarP(&up.resourcePath, "resource", "s", up.resourcePath, "The cluster resource to use (namespace/kind/name format).")
	flags.StringVar(&up.containerName, "container", up.containerName, "The container name to use for remote development")

	_ = command.MarkFlagDirname("local-sync-path")

	flags.Var(
		enumflag.New(&up.syncMode, "sync-mode", SyncModeIds, enumflag.EnumCaseSensitive),
		"sync-mode",
		"Mutagen sync mode.\n"+
			fmt.Sprintf("Available sync modes: %s\n", strings.Join(SyncModeList, ", "))+
			fmt.Sprintf(`"%s" sync mode disables mutagen.`, string(remoteDevMutagenConfig.None)),
	)

	_ = command.RegisterFlagCompletionFunc("sync-mode", cobra.FixedCompletions(SyncModeList, cobra.ShellCompDirectiveDefault))

	up.addContainerConfigFlags(flags)
	up.manager.UpdateFlagSet(command, flags)

	command.MarkFlagsRequiredTogether("container", "rdev-profile")
}

func (up *Options) addContainerConfigFlags(
	flags *pflag.FlagSet,
) {
	flags.StringArrayVar(&up.environPairs, "env", up.environPairs, "Environment variables to set in the remote development container")

	flags.StringVar(&up.limitCPU, "limit-cpu", up.limitCPU, "CPU resource limits for the remote development container")
	flags.StringVar(&up.limitMemory, "limit-memory", up.limitMemory, "Memory resource limits for the remote development container")
	flags.StringVar(&up.requestCPU, "request-cpu", up.requestCPU, "CPU resource reservations for the remote development container")
	flags.StringVar(&up.requestMemory, "request-memory", up.requestMemory, "Memory resource reservations for the remote development container")

	flags.StringSliceVarP(
		&up.portMappings,
		"port-forward",
		"f",
		up.portMappings,
		"Port forward: '8080>3000'\nReverse port forward: '9003<9003'\nComma separated: '8080>3000,9003<9003'",
	)
}
