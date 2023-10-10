package config

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
)

func (manager *Manager) readConfig(fileName string) error {
	manager.viper.SetConfigFile(fileName)

	if err := manager.viper.ReadInConfig(); err != nil {
		// @review update go:1.20 errors.join
		return fmt.Errorf("%w: %s", ErrConfigLoad, err.Error())
	}

	if err := manager.viper.Unmarshal(manager.config); err != nil {
		// @review update go:1.20 errors.join
		return fmt.Errorf("%w: %s", ErrConfigLoad, err.Error())
	}

	for name, profile := range manager.config.Profiles {
		profile.Name = name

		manager.config.Profiles[name] = profile
	}

	return nil
}

func (manager *Manager) importEnvOnly() {
	manager.options.Timeout.ValueOr(func(flag *pflag.Flag) time.Duration {
		if manager.viper.IsSet(flag.Name) {
			return manager.viper.GetDuration(flag.Name)
		}

		return 0
	})
}

func (manager *Manager) importConfig(config *Config) {
	manager.options.OutputFormat.ValueOr(func(flag *pflag.Flag) string {
		if manager.viper.IsSet(flag.Name) {
			return manager.viper.GetString(flag.Name)
		}

		return config.OutputFormat
	})
	manager.options.ProfileName.ValueOr(func(flag *pflag.Flag) string {
		if manager.viper.IsSet(flag.Name) {
			return manager.viper.GetString(flag.Name)
		}

		return config.DefaultProfile
	})
	manager.options.Debug.ValueOr(func(flag *pflag.Flag) bool {
		if manager.viper.IsSet(flag.Name) {
			return manager.viper.GetBool(flag.Name)
		}

		return config.Debug
	})

	if manager.settings.Profile.Name == "" {
		manager.importProfile(&Profile{})

		return
	}

	profile, err := manager.config.getProfile(manager.settings.Profile.Name)
	if err != nil {
		manager.Error = err
		profile = &Profile{}
	}

	manager.importProfile(profile)
}

func (manager *Manager) importProfile(profile *Profile) {
	manager.options.Host.ValueOr(func(flag *pflag.Flag) string {
		if manager.viper.IsSet(flag.Name) {
			return manager.viper.GetString(flag.Name)
		}

		return profile.Host
	})
	manager.options.Scheme.ValueOr(func(flag *pflag.Flag) string {
		if manager.viper.IsSet(flag.Name) {
			return manager.viper.GetString(flag.Name)
		}

		return profile.Scheme
	})
	manager.options.Token.ValueOr(func(flag *pflag.Flag) string {
		if manager.viper.IsSet(flag.Name) {
			return manager.viper.GetString(flag.Name)
		}

		return profile.Token
	})

	manager.importContext(profile.Context)
}

func (manager *Manager) importContext(context Context) {
	manager.options.Organization.ValueOr(func(flag *pflag.Flag) string {
		if manager.viper.IsSet(flag.Name) {
			return manager.viper.GetString(flag.Name)
		}

		return context.Organization
	})
	manager.options.Project.ValueOr(func(flag *pflag.Flag) string {
		if manager.viper.IsSet(flag.Name) {
			return manager.viper.GetString(flag.Name)
		}

		return context.Project
	})
	manager.options.Environment.ValueOr(func(flag *pflag.Flag) string {
		if manager.viper.IsSet(flag.Name) {
			return manager.viper.GetString(flag.Name)
		}

		return context.Environment
	})
	manager.options.ServiceComponent.ValueOr(func(flag *pflag.Flag) string {
		if manager.viper.IsSet(flag.Name) {
			return manager.viper.GetString(flag.Name)
		}

		return context.ServiceComponent
	})
}
