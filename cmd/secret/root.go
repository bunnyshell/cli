package secret

import (
	"os"

	"bunnyshell.com/cli/pkg/config"
	"github.com/spf13/cobra"
)

var mainCmd = &cobra.Command{
	Use:     "secrets",
	Aliases: []string{"sec"},

	Short: "Secrets",
	Long:  "Bunnyshell Secrets",
}

func init() {
	config.MainManager.CommandWithAPI(mainCmd)
}

func GetMainCommand() *cobra.Command {
	return mainCmd
}

func isStdinPresent() (bool, error) {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false, err
	}

	return (fi.Mode() & os.ModeCharDevice) == 0, nil
}
