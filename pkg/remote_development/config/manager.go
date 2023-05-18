package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Manager struct {
	config *Config

	profileName string

	configDirParam  string
	configFileParam string

	configDir  string
	configFile string
}

func NewManager() *Manager {
	return &Manager{
		config: &Config{},

		configDirParam:  defaultConfigDirParam,
		configFileParam: defaultConfigFileParam,
	}
}

func (manager *Manager) Validate() error {
	if manager.configDirParam == defaultConfigDirParam && manager.configFileParam == defaultConfigFileParam {
		return nil
	}

	return manager.ensureConfigFile()
}

func (manager *Manager) Load() error {
	if err := manager.ensureConfigFile(); err != nil {
		return err
	}

	return manager.readConfig()
}

func (manager *Manager) SetProfileName(profileName string) {
	manager.profileName = profileName
}

func (manager *Manager) HasProfileName() bool {
	return manager.profileName != ""
}

func (manager *Manager) GetProfile() (*Profile, error) {
	return manager.config.getProfile(manager.profileName)
}

func (manager *Manager) readConfig() error {
	yamlFile, err := os.ReadFile(filepath.Join(manager.configDir, manager.configFile))
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(yamlFile, manager.config); err != nil {
		return fmt.Errorf("%w: %s", ErrConfigLoad, err.Error())
	}

	for name, profile := range manager.config.Profiles {
		profile.Name = name

		manager.config.Profiles[name] = profile
	}

	return nil
}

func (manager *Manager) ensureConfigFile() error {
	if manager.configDir != "" || manager.configFile != "" {
		return nil
	}

	if filepath.IsAbs(manager.configFileParam) {
		return manager.configToOSPath(manager.configFileParam)
	}

	if err := manager.discoverConfigDir(); err != nil {
		return err
	}

	return manager.configToOSPath(filepath.Join(manager.configDirParam, manager.configFileParam))
}

func (manager *Manager) discoverConfigDir() error {
	if !strings.HasPrefix(manager.configDirParam, ".../") {
		return nil
	}

	subdir := manager.configDirParam[4:]

	workingDirectory, err := os.Getwd()
	if err != nil {
		return err
	}

	for {
		configFileDir := filepath.Join(workingDirectory, subdir)
		configFilePath := filepath.Join(configFileDir, manager.configFileParam)

		if _, err = os.Stat(configFilePath); err == nil {
			manager.configDirParam = configFileDir

			return nil
		}

		parentDirectory := filepath.Dir(workingDirectory)

		if parentDirectory == workingDirectory {
			return fmt.Errorf("%w: unable to find %s", ErrConfigLoad, manager.configDirParam)
		}

		workingDirectory = parentDirectory
	}
}

func (manager *Manager) configToOSPath(file string) error {
	manager.configDir = filepath.Dir(file)
	manager.configFile = filepath.Base(file)

	_, err := os.Stat(filepath.Join(manager.configDir, manager.configFile))

	return err
}

func (manager *Manager) MakeAbsolute(path *string) error {
	if path == nil || *path == "" {
		return nil
	}

	if filepath.IsAbs(*path) {
		return nil
	}

	if (*path)[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		*path = home + (*path)[1:]

		return nil
	}

	*path = filepath.Join(manager.configDir, *path)

	return nil
}
