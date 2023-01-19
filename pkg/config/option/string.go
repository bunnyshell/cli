package option

import (
	"github.com/spf13/pflag"
)

type String struct {
	Group

	value *Value
}

type StringGenerator func(flag *pflag.Flag) string

func NewStringOption(val *string) *String {
	return NewStringValueOption(newStringValue(val))
}

func NewStringValueOption(value Value) *String {
	return &String{
		Group: Group{},

		value: &value,
	}
}

func (option *String) Var() *Value {
	return option.value
}

// @review needs an error variant
func (option *String) ValueOr(generator StringGenerator) string {
	if option.IsChanged() {
		return option.value.String()
	}

	flag := option.GetMainFlag()
	value := generator(flag)

	if value == "" {
		return option.value.String()
	}

	if err := option.updateFlags(value); err != nil {
		// @review not elegant
		panic(err)
	}

	return value
}

func (option *String) AddFlag(name string, usage string) *pflag.Flag {
	return option.AddFlagShort(name, "", usage)
}

func (option *String) AddFlagShort(name string, short string, usage string) *pflag.Flag {
	flag := option.makeFlag(name, short, usage)

	option.Group.AddFlag(flag)

	return flag
}

func (option *String) CloneMainFlag() *pflag.Flag {
	if option.main == nil {
		return nil
	}

	return option.AddFlagShort(option.main.Name, option.main.Shorthand, option.main.Shorthand)
}

func (option *String) makeFlag(name string, short string, usage string) *pflag.Flag {
	return &pflag.Flag{
		Name:      name,
		Shorthand: short,
		Usage:     usage,

		Value:    *option.value,
		DefValue: option.value.String(),
	}
}
