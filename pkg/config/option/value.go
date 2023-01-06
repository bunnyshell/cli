package option

import (
	"github.com/spf13/pflag"
)

type Validator func(string, pflag.Value) error

type Value struct {
	pflag pflag.Value

	Validator Validator
}

func (v Value) String() string {
	return v.pflag.String()
}

func (v Value) Type() string {
	return v.pflag.Type()
}

func (v Value) Set(data string) error {
	if v.Validator != nil {
		if err := v.Validator(data, v.pflag); err != nil {
			return err
		}
	}

	return v.pflag.Set(data)
}
