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

		configDirParam:  ".../.bunnyshell",
		configFileParam: "rdev.yaml",
	}
}

func (manager *Manager) Load() error {
	if err := manager.ensureConfigDir(); err != nil {
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
	yamlFile, err := os.ReadFile(manager.configDir + "/" + manager.configFile)
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

func (manager *Manager) ensureConfigDir() error {
	if filepath.IsAbs(manager.configFileParam) {
		manager.configToOSPath(manager.configFileParam)

		return nil
	}

	if err := manager.discoverConfigDir(); err != nil {
		return err
	}

	manager.configToOSPath(manager.configDirParam + "/" + manager.configFileParam)

	return nil
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
		configFileDir := workingDirectory + "/" + subdir
		configFilePath := configFileDir + "/" + manager.configFileParam

		if _, err = os.Stat(configFilePath); err == nil {
			manager.configDirParam = configFileDir

			return nil
		}

		if workingDirectory == "/" {
			return fmt.Errorf("%w: unable to find %s", ErrConfigLoad, manager.configDirParam)
		}

		workingDirectory = filepath.Dir(workingDirectory)
	}
}

func (manager *Manager) configToOSPath(file string) {
	manager.configDir = filepath.Dir(file)
	manager.configFile = filepath.Base(file)
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

	*path = manager.configDir + "/" + *path

	return nil
}
