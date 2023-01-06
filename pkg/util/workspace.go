package util

import (
	"os"
	"path/filepath"
)

const workspaceDirname = ".bunnyshell"

func GetWorkspaceDirAndShort() (string, string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", "", err
	}

	if home == "/" {
		return "/bunnyshell", "/bunnyshell", nil
	}

	// @review $HOME is more linuxy check os.UserHomeDir and update ?
	return filepath.Join(home, workspaceDirname), filepath.Join("$HOME", workspaceDirname), nil
}

func GetWorkspaceDir() (string, error) {
	home, _, err := GetWorkspaceDirAndShort()

	return home, err
}
