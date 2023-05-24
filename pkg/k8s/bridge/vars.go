package bridge

import "errors"

var (
	ErrNotLoaded = errors.New("bunnyshell not loaded")

	ErrResourceNotFound = errors.New("resource %s not found")

	ErrNoComponentResources = errors.New("no component resources available")
)
