package template

import (
	"path/filepath"
)

type Data struct {
	Directory string
}

func (d Data) Basename() string {
	return filepath.Base(d.Directory)
}
