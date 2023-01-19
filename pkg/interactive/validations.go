package interactive

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

func All(funcs ...survey.Validator) survey.Validator {
	return func(input interface{}) error {
		for _, callable := range funcs {
			err := callable(input)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

func Lowercase() survey.Validator {
	return func(input interface{}) error {
		str, ok := input.(string)
		if !ok {
			return ErrInvalidValue
		}

		if strings.ToLower(str) != input {
			return fmt.Errorf("%w: must be lowercase", ErrInvalidValue)
		}

		return nil
	}
}

func AssertBetween(min int32, max int32) survey.Validator {
	return func(input interface{}) error {
		str, ok := input.(string)
		if !ok {
			return ErrInvalidValue
		}

		i, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			return fmt.Errorf("%w: must be an integer", ErrInvalidValue)
		}

		val := int32(i)
		if val < min || val > max {
			return fmt.Errorf("%w: must be between %d and %d", ErrInvalidValue, min, max)
		}

		return nil
	}
}

func AssertMinimumLength(length int) survey.Validator {
	return survey.MinLength(length)
}
