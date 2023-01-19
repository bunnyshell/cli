package option

import (
	"time"

	"github.com/spf13/pflag"
)

type Duration struct {
	String
}

type DurationGenerator func(flag *pflag.Flag) time.Duration

func NewDurationOption(val *time.Duration) *Duration {
	value := newDurationValue(val)

	return &Duration{
		String: *NewStringValueOption(value),
	}
}

func (option *Duration) ValueOr(generator DurationGenerator) string {
	return option.String.ValueOr(func(flag *pflag.Flag) string {
		value := generator(flag).String()

		if value == "0s" {
			return ""
		}

		return value
	})
}
