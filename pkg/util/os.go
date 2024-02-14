package util

import "os"

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func IsStdinPresent() (bool, error) {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false, err
	}

	return (fi.Mode() & os.ModeCharDevice) == 0, nil
}
