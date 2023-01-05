package util

import (
	"os"
	"path/filepath"
)

const workspaceDirname = ".bunnyshell"

func GetWorkspaceDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	if home == "/" {
		return "/bunnyshell", nil
	}

	return filepath.Join(home, workspaceDirname), nil
}
