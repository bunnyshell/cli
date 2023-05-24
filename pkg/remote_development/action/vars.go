package action

import "errors"

var (
	ErrOneSyncPathSupported = errors.New("only one sync path is supported")

	ErrRemoteDevNotInitialized = errors.New("call Up.Run() successfully before calling other methods")
)
