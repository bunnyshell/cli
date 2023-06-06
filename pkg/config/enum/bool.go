package enum

import (
	"github.com/spf13/pflag"
	"github.com/thediveo/enumflag/v2"
)

type Bool enumflag.Flag

const (
	BoolNone Bool = iota
	BoolTrue
	BoolFalse
)

var BoolMap = map[Bool][]string{
	BoolNone:  {"inherit"},
	BoolTrue:  {"true"},
	BoolFalse: {"false"},
}

var BoolList = []string{
	"true",
	"false",
}

func NewBoolValue(value *Bool) pflag.Value {
	return enumflag.New(value, "bool", BoolMap, enumflag.EnumCaseInsensitive)
}

func BoolFlag(value *Bool, name string, usage string) *pflag.Flag {
	enumValue := NewBoolValue(value)

	return &pflag.Flag{
		Name:      name,
		Shorthand: "",
		Usage:     usage,
		Value:     enumValue,
		DefValue:  enumValue.String(),
	}
}
