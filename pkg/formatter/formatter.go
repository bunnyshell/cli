package formatter

import (
	"encoding/json"
	"errors"
	"fmt"

	"bunnyshell.com/sdk"
	"gopkg.in/yaml.v3"
)

var errUnknownFormat = errors.New("unknown format")

func Formatter(data interface{}, format string) ([]byte, error) {
	if format == "stylish" {
		return stylish(data)
	}

	// not really useful
	// maybe update the hal spec ?
	switch t := data.(type) {
	case *sdk.PaginatedOrganizationCollection:
		t.Links = nil
	case *sdk.PaginatedProjectCollection:
		t.Links = nil
	case *sdk.PaginatedEnvironmentCollection:
		t.Links = nil
	case *sdk.PaginatedComponentCollection:
		t.Links = nil
	case *sdk.PaginatedEventCollection:
		t.Links = nil
	case *sdk.PaginatedEnvironmentVariableCollection:
		t.Links = nil
	}

	switch format {
	case "json":
		return JSONFormatter(data)
	case "yaml", "yml":
		return YAMLFormatter(data)
	}

	return nil, fmt.Errorf("%w: %s", errUnknownFormat, format)
}

func JSONFormatter(data interface{}) ([]byte, error) {
	return json.MarshalIndent(data, "", "  ")
}

func YAMLFormatter(data interface{}) ([]byte, error) {
	return yaml.Marshal(data)
}
