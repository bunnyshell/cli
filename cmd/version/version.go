package version

import (
	"errors"
	"net/http"
	"strings"

	"github.com/spf13/cobra"

	"bunnyshell.com/cli/pkg/build"
	"bunnyshell.com/cli/pkg/lib"
)

var ClientOnly = false

var mainCmd = &cobra.Command{
	Use:   "version",
	Short: "Version Information",
	RunE: func(cmd *cobra.Command, args []string) error {
		data := map[string]interface{}{}

		data["client"] = map[string]interface{}{
			"version": build.Version,
		}

		if !ClientOnly {
			latestRelease, err := getLatestRelease()
			if err != nil {
				return err
			}

			if latestRelease == build.Version {
				cmd.Println("You are using the latest version: " + latestRelease)
				return nil
			}

			data["server"] = map[string]interface{}{
				"version": latestRelease,
			}
		}

		return lib.FormatCommandData(cmd, data)
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
