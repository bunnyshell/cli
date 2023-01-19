package option

import (
	"strconv"

	"github.com/spf13/pflag"
)

type Count struct {
	String
}

type CountGenerator func(flag *pflag.Flag) int

func NewCountOption(val *int) *Count {
	value := newCountValue(val)

	return &Count{
		String: *NewStringValueOption(value),
	}
}
func (option *Count) AddFlag(name string, usage string) *pflag.Flag {
	return option.AddFlagShort(name, "", usage)
}

func (option *Count) AddFlagShort(name string, short string, usage string) *pflag.Flag {
	flag := option.String.AddFlagShort(name, short, usage)

	flag.NoOptDefVal = "+1"

	return flag
}

func (option *Count) ValueOr(generator CountGenerator) string {
	return option.String.ValueOr(func(flag *pflag.Flag) string {
		return strconv.Itoa(generator(flag))
	})
}
