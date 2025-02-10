package up

import (
	"errors"
	"time"
)

const defaultWaitTimeout = 120 * time.Second

var (
	ErrResourceLoaderNotHydrated = errors.New("resourceLoader needs to be hydrated")
)
