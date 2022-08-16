package lib

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Profile struct {
	Host    string  `json:"host,omitempty" yaml:"host,omitempty"`
	Token   string  `json:"token,omitempty" yaml:"token,omitempty"`
	Context Context `json:"context,omitempty" yaml:"context,omitempty"`
}

type NamedProfiles map[string]Profile

type Context struct {
	Organization     string `json:"organization,omitempty" yaml:"organization,omitempty"`
	Project          string `json:"project,omitempty" yaml:"project,omitempty"`
	Environment      string `json:"environment,omitempty" yaml:"environment,omitempty"`
	ServiceComponent string `json:"serviceComponent,omitempty" yaml:"serviceComponent,omitempty"`
}

type Config struct {
	Debug          bool          `json:"debug" yaml:"debug"`
	OutputFormat   string        `json:"outputFormat,omitempty" yaml:"outputFormat,omitempty" binding:"required"`
	DefaultProfile string        `json:"defaultProfile" yaml:"defaultProfile" binding:"required"`
	Timeout        time.Duration `json:",omitempty" yaml:",omitempty"`
	Profiles       NamedProfiles `json:",omitempty" yaml:",omitempty"`
}

func GetConfig() (*Config, error) {
	var config Config
	err := viper.Unmarshal(&config)

	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal: %s", err)
	}

	return &config, nil
}

func SetDefaultProfile(name string) error {
	config, err := GetConfig()
	if err != nil {
		return err
	}

	if _, ok := config.Profiles[name]; !ok {
		return fmt.Errorf("unknown profile %s", name)
	}

	viper.Set("defaultProfile", name)

	return saveConfig()
}

func NewProfile() *Profile {
	return &Profile{}
}

func GetProfile(name string) (*Profile, error) {
	config, err := GetConfig()
	if err != nil {
		return nil, err
	}

	profile, ok := config.Profiles[name]
	if !ok {
		return nil, errors.New("profile not found")
	}

	return &profile, nil
}

func AddProfile(profile Profile, name string) error {
	config, err := GetConfig()
	if err != nil {
		return err
	}

	if _, ok := config.Profiles[name]; ok {
		return fmt.Errorf("%s already exists", name)
	}

	if config.Profiles == nil {
		config.Profiles = NamedProfiles{}
	}

	config.Profiles[name] = profile
	viper.Set("profiles", config.Profiles)

	return saveConfig()
}

func RemoveProfile(name string) error {
	config, err := GetConfig()
	if err != nil {
		return err
	}

	delete(config.Profiles, name)
	viper.Set("profiles", config.Profiles)

	if config.DefaultProfile == name {
		config.DefaultProfile = ""
		viper.Set("defaultProfile", config.DefaultProfile)
	}

	return saveConfig()
}

func saveConfig() error {
	return viper.WriteConfig()
}
