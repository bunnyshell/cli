package action

import "errors"

var (
	ErrDebugCmpNotInitialized = errors.New("call Up.Run() successfully before calling other methods")
)
