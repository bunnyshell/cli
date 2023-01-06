package option

import (
	"github.com/spf13/pflag"
)

type Bool struct {
	String
}

const trueStr = "true"

type BoolGenerator func(flag *pflag.Flag) bool

func NewBoolOption(val *bool) *Bool {
	value := newBoolValue(val)

	return &Bool{
		String: *NewStringValueOption(value),
	}
}

func (option *Bool) AddFlag(name string, usage string) *pflag.Flag {
	return option.AddFlagShort(name, "", usage)
}

func (option *Bool) AddFlagShort(name string, short string, usage string) *pflag.Flag {
	flag := option.String.AddFlagShort(name, short, usage)

	flag.NoOptDefVal = trueStr

	return flag
}

func (option *Bool) ValueOr(generator BoolGenerator) string {
	return option.String.ValueOr(func(flag *pflag.Flag) string {
		value := generator(flag)

		if !value {
			return ""
		}

		return trueStr
	})
}
