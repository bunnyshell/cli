package formatter

import (
	"encoding/json"
	"errors"

	"bunnyshell.com/sdk"
	"gopkg.in/yaml.v3"
)

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
	}

	switch format {
	case "json":
		return JsonFormatter(data)
	case "yaml", "yml":
		return YamlFormatter(data)
	}

	return nil, errors.New("Unknown format: " + format)
}

func JsonFormatter(data interface{}) ([]byte, error) {
	return json.MarshalIndent(data, "", "  ")
}

func YamlFormatter(data interface{}) ([]byte, error) {
	return yaml.Marshal(data)
}
