package ssh

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/kevinburke/ssh_config"
	"golang.org/x/crypto/ssh"
)

func GetConfig() (*ssh_config.Config, error) {
	configFile, err := getConfigFile()
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return &ssh_config.Config{}, nil
	}

	if err != nil {
		return nil, err
	}

	return ssh_config.Decode(configFile)
}

func SaveConfig(cfg *ssh_config.Config) error {
	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	data, err := cfg.MarshalText()
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

func RemoveHost(cfg *ssh_config.Config, hostname string) error {
	newHosts := []*ssh_config.Host{}
	for _, host := range cfg.Hosts {
		if host.Matches(hostname) {
			continue
		}

		newHosts = append(newHosts, host)
	}

	cfg.Hosts = newHosts

	return nil
}

func NewKV(paramName, paramValue string) *ssh_config.KV {
	return &ssh_config.KV{
		Key:   "  " + paramName,
		Value: paramValue,
	}
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ".ssh", "config"), nil
}

func getConfigFile() (*os.File, error) {
	filePath, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}

	return os.Open(filePath)
}

func PrivateKeyFile(file string) ssh.AuthMethod {
	buffer, err := os.ReadFile(file)
	if err != nil {
		return nil
	}
	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}
