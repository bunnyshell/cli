package up

import (
	"errors"
	"time"
)

const defaultWaitTimeout = 120 * time.Second

var (
	ErrNoRDevConfig = errors.New("component does not have a remote development config")

	ErrEmptyList = errors.New("empty list")

	ErrOneSyncPathSupported = errors.New("only one sync path is supported")

	ErrTooManySimpleConfig = errors.New("too many simple resources")

	ErrTooManySimpleResources = errors.New("too many simple containers")

	ErrUnknownProfileType       = errors.New("unknown profile type")
	ErrUnknownConfigurationType = errors.New("unknown configuration type")

	ErrResourceLoaderNotHydrated = errors.New("resourceLoader needs to be hydrated")
)
