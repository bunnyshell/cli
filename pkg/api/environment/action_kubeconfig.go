package environment

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/api/common"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
)

type KubeConfigData map[string]any

type KubeConfigItem struct {
	Data *sdk.EnvironmentKubeConfigKubeConfigRead

	Bytes []byte
}

type KubeConfigOptions struct {
	common.ItemOptions

	OverrideClusterServer string
}

func NewKubeConfigOptions(id string, overrideClusterServer string) *KubeConfigOptions {
	return &KubeConfigOptions{
		ItemOptions:           *common.NewItemOptions(id),
		OverrideClusterServer: overrideClusterServer,
	}
}

func KubeConfig(options *KubeConfigOptions) (*KubeConfigItem, error) {
	model, resp, err := KubeConfigRaw(&options.ItemOptions)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	var cfgBytes []byte
	if options.OverrideClusterServer != "" {
		// normally it would be:
		//     model.Clusters[0].Cluster.Server, err = overrideServer(model.Clusters[0].Cluster.Server, options.OverrideClusterServer)
		//     if err != nil {
		//			return nil, fmt.Errorf("error overriding cluster server: %w", err)
		//		}
		//     cfgBytes, err = yaml.Marshal(model)
		//     if err != nil {
		//         return nil, fmt.Errorf("error marshaling to YAML: %w", err)
		//     }
		// but model is not unmarshalled correctly, because the struct doesn't have yaml tags, only json tags,
		// and the response is in application/x+yaml format (it skips the multi words properties like ApiVersion, CurrentContext)
		// se we use the string replace hack

		tmpBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		oldServer := model.Clusters[0].Cluster.Server
		model.Clusters[0].Cluster.Server, err = overrideServer(model.Clusters[0].Cluster.Server, options.OverrideClusterServer)
		if err != nil {
			return nil, fmt.Errorf("error overriding cluster server: %w", err)
		}

		cfgBytes = bytes.ReplaceAll(tmpBytes, []byte(oldServer), []byte(model.Clusters[0].Cluster.Server))
	} else {
		cfgBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
	}

	return &KubeConfigItem{
		Data:  model,
		Bytes: cfgBytes,
	}, nil
}

func KubeConfigRaw(options *common.ItemOptions) (*sdk.EnvironmentKubeConfigKubeConfigRead, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).EnvironmentAPI.EnvironmentKubeConfig(ctx, options.ID)

	return request.Execute()
}

// overrideServer takes a base Kubernetes server URL like
//
//	"https://my.example.com:6443"
//
// and an override string which can be one of:
//
//	"SCHEME://..."    → full replacement (any Scheme, host, port, path, etc.)
//	":PORT"           → change only the port
//	"HOST:PORT"       → change host and port
//	"HOST"       	  → change only the host
//
// It returns the new URL string or an error.
func overrideServer(base, override string) (string, error) {
	// Full-URL override
	if strings.Contains(override, "://") {
		// we treat anything containing "://" as a complete URL replacement
		v, err := url.Parse(override)
		if err != nil {
			return "", fmt.Errorf("parsing full override %q: %w", override, err)
		}
		return v.String(), nil
	}

	// Parse the base
	u, err := url.Parse(base)
	if err != nil {
		return "", fmt.Errorf("parsing base %q: %w", base, err)
	}

	// Extract the existing host and port
	host := u.Hostname()
	port := u.Port()

	switch {
	case strings.HasPrefix(override, ":"):
		// override only the port
		port = strings.TrimPrefix(override, ":")

	case strings.Contains(override, ":"):
		// override both host and port
		h, p, err := net.SplitHostPort(override)
		if err != nil {
			return "", fmt.Errorf("parsing host:port override %q: %w", override, err)
		}
		host, port = h, p

	default:
		// (optional) override only the host, keep existing port
		host = override
	}

	// Re-join host and port (adds the brackets for IPv6 if needed)
	u.Host = net.JoinHostPort(host, port)

	return u.String(), nil
}
