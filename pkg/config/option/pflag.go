package option

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
)

var (
	flagCount = 0

	flagSet = pflag.NewFlagSet("export_pflag_value", pflag.ContinueOnError)
)

// unexported newStringValue()
// @see https://github.com/spf13/pflag/blob/v1.0.5/string.go
func newStringValue(val *string) Value {
	name := getNewFlagName()

	flagSet.StringVar(val, name, *val, "")

	return valueFrom(name)
}

// unexported newBoolValue()
// @see https://github.com/spf13/pflag/blob/v1.0.5/bool.go
func newBoolValue(val *bool) Value {
	name := getNewFlagName()

	flagSet.BoolVar(val, name, *val, "")

	return valueFrom(name)
}

// unexported newCountValue()
// @see https://github.com/spf13/pflag/blob/v1.0.5/count.go
func newCountValue(val *int) Value {
	name := getNewFlagName()

	flagSet.CountVar(val, name, "")

	return valueFrom(name)
}

// unexported newDurationValue()
// @see https://github.com/spf13/pflag/blob/v1.0.5/duration.go
func newDurationValue(val *time.Duration) Value {
	name := getNewFlagName()

	flagSet.DurationVar(val, name, *val, "")

	return valueFrom(name)
}

func getNewFlagName() string {
	flagCount++

	return fmt.Sprintf("flag%d", flagCount)
}

func valueFrom(name string) Value {
	return Value{
		pflag: flagSet.Lookup(name).Value,
	}
}
