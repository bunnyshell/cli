package up

import (
	remoteDevMutagenConfig "bunnyshell.com/dev/pkg/mutagen/config"
	"github.com/thediveo/enumflag/v2"
)

type SyncMode enumflag.Flag

const (
	None SyncMode = iota
	TwoWaySafe
	TwoWayResolved
	OneWaySafe
	OneWayReplica
)

var SyncModeToMutagenMode = map[SyncMode]remoteDevMutagenConfig.Mode{
	None:           remoteDevMutagenConfig.None,
	TwoWaySafe:     remoteDevMutagenConfig.TwoWaySafe,
	TwoWayResolved: remoteDevMutagenConfig.TwoWayResolved,
	OneWaySafe:     remoteDevMutagenConfig.OneWaySafe,
	OneWayReplica:  remoteDevMutagenConfig.OneWayReplica,
}

var SyncModeIds = map[SyncMode][]string{
	None:           {string(remoteDevMutagenConfig.None)},
	TwoWaySafe:     {string(remoteDevMutagenConfig.TwoWaySafe)},
	TwoWayResolved: {string(remoteDevMutagenConfig.TwoWayResolved)},
	OneWaySafe:     {string(remoteDevMutagenConfig.OneWaySafe)},
	OneWayReplica:  {string(remoteDevMutagenConfig.OneWayReplica)},
}

var SyncModeList = []string{
	string(remoteDevMutagenConfig.None),
	string(remoteDevMutagenConfig.TwoWaySafe),
	string(remoteDevMutagenConfig.TwoWayResolved),
	string(remoteDevMutagenConfig.OneWaySafe),
	string(remoteDevMutagenConfig.OneWayReplica),
}
