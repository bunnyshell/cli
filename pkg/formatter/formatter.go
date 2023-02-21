package formatter

import (
	"encoding/json"
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

var errUnknownFormat = errors.New("unknown format")

func Formatter(data interface{}, format string) ([]byte, error) {
	if format == "stylish" {
		return stylish(data)
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
