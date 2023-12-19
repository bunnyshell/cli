package git

import (
	"net/url"
	"strings"
)

func ParseGitSec(spec string) (string, string, error) {
	if spec[0] == '@' {
		return "", spec[1:], nil
	}

	info, err := url.Parse(spec)
	if err != nil {
		return "", "", err
	}

	if !strings.Contains(info.Path, "@") {
		return spec, "", nil
	}

	chunks := strings.SplitN(info.Path, "@", 2)
	info.Path = chunks[0]

	return info.String(), chunks[1], nil
}
