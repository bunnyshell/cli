package remote_development

import (
	"io"
	"os"
	"path/filepath"

	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/util"
)

const (
	KubeConfigFilename = "kube-config.yaml"
)

func (r *RemoteDevelopment) ensureEnvironmentWorkspaceDir() error {
	workspace, err := util.GetWorkspaceDir()
	if err != nil {
		return err
	}

	r.WithEnvironmentWorkspaceDir(filepath.Join(workspace, r.environment.GetId()))
	return os.MkdirAll(r.environmentWorkspaceDir, 0755)
}

func (r *RemoteDevelopment) ensureEnvironmentKubeConfig() error {
	kubeConfigPath := filepath.Join(r.environmentWorkspaceDir, KubeConfigFilename)
	if err := downloadEnvironmentKubeConfig(kubeConfigPath, r.environment.GetId()); err != nil {
		return err
	}
	r.WithKubeConfigPath(kubeConfigPath)

	return nil
}

func downloadEnvironmentKubeConfig(kubeConfigPath, environmentId string) error {
	kubeConfigFile, err := os.Create(kubeConfigPath)
	if err != nil {
		return err
	}
	defer kubeConfigFile.Close()

	ctx, cancel := lib.GetContext()
	defer cancel()

	request := lib.GetAPI().EnvironmentApi.EnvironmentKubeConfig(ctx, environmentId)
	_, resp, err := request.Execute()
	if err != nil && err.Error() != "undefined response type" {
		return err
	}

	_, err = io.Copy(kubeConfigFile, resp.Body)
	return err
}
