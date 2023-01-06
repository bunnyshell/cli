package lib

import (
	"context"

	"bunnyshell.com/cli/pkg/build"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/net"
	"bunnyshell.com/sdk"
)

func GetProfile() config.Profile {
	return config.GetSettings().Profile
}

func GetAPI() *sdk.APIClient {
	return GetAPIFromProfile(GetProfile())
}

func GetContext() (context.Context, context.CancelFunc) {
	return GetContextFromProfile(GetProfile())
}

func GetAPIFromProfile(profile config.Profile) *sdk.APIClient {
	return sdk.NewAPIClient(profileToConfiguration(profile))
}

func GetContextFromProfile(profile config.Profile) (context.Context, context.CancelFunc) {
	ctx := context.WithValue(context.Background(), sdk.ContextAPIKeys, map[string]sdk.APIKey{
		"ApiKeyAuth": {
			Key: profile.Token,
		},
	})

	timeout := config.GetSettings().Timeout
	if timeout == 0 {
		return ctx, func() {}
	}

	return context.WithTimeout(ctx, timeout)
}

func profileToConfiguration(profile config.Profile) *sdk.Configuration {
	configuration := getDefaultConfiguration()

	if profile.Host != "" {
		configuration.Host = profile.Host
	}

	return configuration
}

func getDefaultConfiguration() *sdk.Configuration {
	configuration := sdk.NewConfiguration()

	configuration.UserAgent = "BunnyCLI+" + build.Version + "/" + configuration.UserAgent
	configuration.Debug = config.GetSettings().Debug
	configuration.HTTPClient = net.GetCLIClient()

	return configuration
}
