package ssh

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"bunnyshell.com/cli/pkg/util"
	"github.com/kevinburke/ssh_config"
	"golang.org/x/crypto/ssh"
)

func GetConfig() (*ssh_config.Config, error) {
	configFile, err := getBunnyshellConfigFile()
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return &ssh_config.Config{}, nil
	}
	defer configFile.Close()

	if err != nil {
		return nil, err
	}

	return ssh_config.Decode(configFile)
}

func SaveConfig(cfg *ssh_config.Config) error {
	filePath, err := getBunnyshellConfigFilePath()
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

func IncludeBunnyshellConfig() error {
	bunnyshellConfigFilePath, err := getBunnyshellConfigFilePath()
	if err != nil {
		return err
	}

	includeDirective := fmt.Sprintf("Include %s", bunnyshellConfigFilePath)
	isIncluded, err := isBunnyshellConfigIncluded(includeDirective)
	if err != nil {
		return err
	}
	if isIncluded {
		return nil
	}
	includeBunnyshellConfig(includeDirective)

	return nil
}

func isBunnyshellConfigIncluded(includeDirective string) (bool, error) {
	file, err := getConfigFile()
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() == includeDirective {
			return true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return false, nil
}

func includeBunnyshellConfig(includeDirective string) error {
	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	appendData := fmt.Sprintf("\n%s\n%s\n", "# do not edit - generated by bunnyshell", includeDirective)
	if _, err := file.WriteString(appendData); err != nil {
		return err
	}

	return nil
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

func getBunnyshellConfigFilePath() (string, error) {
	workspaceDir, err := util.GetWorkspaceDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(workspaceDir, "ssh-config"), nil
}

func getBunnyshellConfigFile() (*os.File, error) {
	filePath, err := getBunnyshellConfigFilePath()
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
