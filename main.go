package main

import (
	"errors"
	"fmt"
	"os"

	"bunnyshell.com/cli/cmd"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/interactive"
)

func main() {
	defer recovery()

	cmd.Execute()
}

func recovery() {
	rErr := recover()
	if rErr == nil {
		return
	}

	if err, ok := rErr.(error); ok {
		if errors.Is(err, interactive.ErrInvalidValue) {
			fmt.Println("[panic]", err)
			os.Exit(1)
		}

		if errors.Is(err, config.ErrInvalidValue) {
			fmt.Println("[panic]", err)
			os.Exit(1)
		}
	}

	panic(rErr)
}
