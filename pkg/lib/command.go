package lib

import (
	"errors"

	"bunnyshell.com/cli/pkg/config"
	"github.com/spf13/cobra"
)

var ErrNotStylish = errors.New("only stylish format is supported")

func OnlyStylish(cmd *cobra.Command, args []string) error {
	if config.GetSettings().IsStylish() {
		return nil
	}

	return ErrNotStylish
}
