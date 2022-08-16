package lib

import (
	"context"

	"bunnyshell.com/cli/pkg/build"
	"bunnyshell.com/cli/pkg/net"
	"bunnyshell.com/sdk"
)

func GetAPI() *sdk.APIClient {
	return GetApiFromProfile(CLIContext.Profile)
}

func GetContext() (context.Context, context.CancelFunc) {
	return GetContextFromProfile(CLIContext.Profile)
}

func GetApiFromProfile(profile Profile) *sdk.APIClient {
	return sdk.NewAPIClient(profileToConfiguration(profile))
}

func GetContextFromProfile(profile Profile) (context.Context, context.CancelFunc) {
	ctx := context.WithValue(context.Background(), sdk.ContextAPIKeys, map[string]sdk.APIKey{
		"ApiKeyAuth": {
			Key: profile.Token,
		},
	})

	if CLIContext.Timeout == 0 {
		return ctx, func() {}
	}

	return context.WithTimeout(ctx, CLIContext.Timeout)
}

func profileToConfiguration(profile Profile) *sdk.Configuration {
	configuration := getDefaultConfiguration()

	if profile.Host != "" {
		configuration.Host = profile.Host
	}

	return configuration
}

func getDefaultConfiguration() *sdk.Configuration {
	configuration := sdk.NewConfiguration()

	configuration.UserAgent = "BunnyCLI+" + build.Version + "/" + configuration.UserAgent
	configuration.Debug = CLIContext.Debug
	configuration.HTTPClient = net.GetCLIClient()

	return configuration
}
