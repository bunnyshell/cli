package util

import "os"

const BunnyshellWorkspaceDirname = ".bunnyshell"

func GetWorkspaceDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	if home == "/" {
		return "/bunnyshell", nil
	}

	return home + "/" + BunnyshellWorkspaceDirname, nil
}
