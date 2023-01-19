package config

import (
	"errors"
	"os"
	"path/filepath"

	"bunnyshell.com/cli/pkg/build"
	"bunnyshell.com/cli/pkg/formatter"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Manager struct {
	Error error

	viper *viper.Viper

	options *Options

	config *Config

	settings *Settings
}

func NewManager() *Manager {
	settings := NewSettings()

	return &Manager{
		viper: viper.New(),

		config: &Config{},

		options: NewOptions(settings),

		settings: settings,
	}
}

func (manager *Manager) SetDefaultProfile(name string) error {
	return manager.config.setDefaultProfile(name)
}

func (manager *Manager) HasProfile(name string) bool {
	_, err := manager.config.getProfile(name)

	return !errors.Is(err, ErrUnknownProfile)
}

func (manager *Manager) GetProfile(name string) (*Profile, error) {
	return manager.config.getProfile(name)
}

func (manager *Manager) SetProfile(profile Profile) {
	manager.config.Profiles[profile.Name] = profile
}

func (manager *Manager) AddProfile(profile Profile) error {
	return manager.config.addProfile(profile)
}

func (manager *Manager) RemoveProfile(name string) error {
	return manager.config.removeProfile(name)
}

func (manager *Manager) Load() {
	manager.viper.SetEnvPrefix(build.EnvPrefix)
	manager.viper.AutomaticEnv()

	configFile := manager.options.ConfigFile.ValueOr(func(flag *pflag.Flag) string {
		return manager.viper.GetString(flag.Name)
	})

	if err := manager.readConfig(configFile); err != nil {
		manager.Error = err
	}

	manager.importEnvOnly()
	manager.importConfig(manager.config)
}

func (manager *Manager) Save() error {
	return manager.save()
}

func (manager *Manager) SafeSave() error {
	exists, err := fileExists(manager.settings.ConfigFile)
	if err != nil {
		return err
	}

	if exists {
		return ErrConfigExists
	}

	return manager.Save()
}

func (manager *Manager) save() error {
	configDir := filepath.Dir(manager.settings.ConfigFile)
	if err := os.MkdirAll(configDir, os.FileMode(configDirPerm)); err != nil {
		return err
	}

	format := getFormatForFile(manager.settings.ConfigFile)

	data, err := formatter.Formatter(manager.config, format)
	if err != nil {
		return err
	}

	return os.WriteFile(manager.settings.ConfigFile, data, os.FileMode(configFilePerm))
}

func getFormatForFile(file string) string {
	ext := filepath.Ext(file)

	switch ext {
	case ".json":
		return "json"
	case ".yaml":
		fallthrough
	default:
		return "yaml"
	}
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}
