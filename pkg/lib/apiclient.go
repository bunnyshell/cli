package lib

import (
	"io"
	"os"
)

func DownloadEnvironmentKubeConfig(kubeConfigPath, environmentID string) error {
	kubeConfigFile, err := os.Create(kubeConfigPath)
	if err != nil {
		return err
	}
	defer kubeConfigFile.Close()

	ctx, cancel := GetContext()
	defer cancel()

	request := GetAPI().EnvironmentApi.EnvironmentKubeConfig(ctx, environmentID)

	_, resp, err := request.Execute()
	if err != nil && err.Error() != "undefined response type" {
		return err
	}

	_, err = io.Copy(kubeConfigFile, resp.Body)

	return err
}
