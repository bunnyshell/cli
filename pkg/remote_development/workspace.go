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

func (r *RemoteDevelopment) EnsureEnvironmentKubeConfig() error {
	workspace, err := util.GetWorkspaceDir()
	if err != nil {
		return err
	}

	kubeConfigPath := filepath.Join(workspace, r.OrganizationId, r.ProjectId, r.EnvironmentId, KubeConfigFilename)
	downloadEnvironmentKubeConfig(kubeConfigPath, r.EnvironmentId)
	r.WithKubernetesClient(kubeConfigPath)

	return nil
}

func (r *RemoteDevelopment) EnsureComponentFolder() error {
	workspace, err := util.GetWorkspaceDir()
	if err != nil {
		return err
	}

	r.ComponentFolderPath = filepath.Join(workspace, r.OrganizationId, r.ProjectId, r.EnvironmentId, r.ComponentName)
	return os.MkdirAll(r.ComponentFolderPath, 0755)
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
