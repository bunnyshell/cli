package down

import (
	"errors"
)

var ErrResourceLoaderNotHydrated = errors.New("resourceLoader needs to be hydrated")
