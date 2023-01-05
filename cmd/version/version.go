package version

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/cobra"

	"bunnyshell.com/cli/pkg/build"
)

var ClientOnly = false

var mainCmd = &cobra.Command{
	Use:   "version",
	Short: "Version Information",

	ValidArgsFunction: cobra.NoFileCompletions,

	Run: func(cmd *cobra.Command, args []string) {
		release := "v" + build.Version
		if ClientOnly {
			cmd.Printf("You are using: %s\n", release)
			return
		}

		latestRelease, err := getLatestRelease()
		if err != nil {
			fmt.Println(err)
			return
		}

		if strings.Contains(release, "-") {
			cmd.Printf("You are running a pre-release version %s. Latest is: %s\n", release, latestRelease)
			return
		}

		if latestRelease == release {
			cmd.Printf("You are using the latest version: %s\n", latestRelease)
			return
		}

		fmt.Printf("Your version %s is older than the latest: %s\n", release, latestRelease)
	},
}

func GetMainCommand() *cobra.Command {
	return mainCmd
}

func init() {
	mainCmd.Flags().BoolVar(&ClientOnly, "client", ClientOnly, "If true, shows client version only (no server required).")
}

func getLatestRelease() (string, error) {
	// Set up the HTTP request
	req, err := http.NewRequest("GET", build.LatestReleaseUrl, nil)
	if err != nil {
		return "", err
	}

	transport := http.Transport{}
	resp, err := transport.RoundTrip(req)
	if err != nil {
		return "", err
	}
	// Check if you received the status codes you expect. There may
	// status codes other than 200 which are acceptable.
	if resp.StatusCode != 200 && resp.StatusCode != 302 {
		return "", errors.New("must be a redirect")
	}

	redirect := resp.Header.Get("Location")
	parts := strings.Split(redirect, "/")
	return parts[len(parts)-1], nil
}
