package version

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"bunnyshell.com/cli/pkg/build"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

var errNotRedirect = errors.New("must be a redirect")

var ClientOnly = false

var mainCmd = &cobra.Command{
	Use:   "version",
	Short: "Version Information",

	ValidArgsFunction: cobra.NoFileCompletions,

	RunE: func(cmd *cobra.Command, args []string) error {
		release := getCurrentRelease()

		if ClientOnly {
			cmd.Printf("You are using: %s\n", release)

			return nil
		}

		latestRelease, err := getLatestRelease()
		if err != nil {
			return lib.FormatCommandError(cmd, err)
		}

		if strings.Contains(release, "-") {
			cmd.Printf("You are running a pre-release version %s. Latest is: %s\n", release, latestRelease)

			return nil
		}

		if latestRelease == release {
			cmd.Printf("You are using the latest version: %s\n", latestRelease)

			return nil
		}

		cmd.Printf("Your version %s is older than the latest: %s\n", release, latestRelease)

		return nil
	},
}

func init() {
	mainCmd.Flags().BoolVar(&ClientOnly, "client", ClientOnly, "If true, shows client version only (no network required)")
}

func GetMainCommand() *cobra.Command {
	return mainCmd
}

func getCurrentRelease() string {
	if build.Version[0] == 'v' {
		return build.Version
	}

	return "v" + build.Version
}

func getLatestRelease() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), lib.CLIContext.Timeout)
	defer cancel()

	// Set up the HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, build.LatestReleaseUrl, nil)
	if err != nil {
		return "", err
	}

	resp, err := (&http.Transport{}).RoundTrip(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound {
		return "", fmt.Errorf("%w (status code: %d)", errNotRedirect, resp.StatusCode)
	}

	redirect := resp.Header.Get("Location")
	parts := strings.Split(redirect, "/")

	return parts[len(parts)-1], nil
}
