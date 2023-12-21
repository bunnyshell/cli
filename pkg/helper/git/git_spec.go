package git

import (
	"errors"
	"net/url"
	"strings"
)

var errEmptySpec = errors.New("empty spec")

func ParseGitSec(spec string) (string, string, error) {
	if len(spec) == 0 {
		return "", "", errEmptySpec
	}

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
