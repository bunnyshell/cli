package remote_development

import (
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

	r.WithEnvironmentWorkspaceDir(filepath.Join(workspace, r.environmentResource.Environment.GetId()))

	return os.MkdirAll(r.environmentWorkspaceDir, 0755)
}

func (r *RemoteDevelopment) ensureEnvironmentKubeConfig() error {
	kubeConfigPath := filepath.Join(r.environmentWorkspaceDir, KubeConfigFilename)
	if err := lib.DownloadEnvironmentKubeConfig(kubeConfigPath, r.environmentResource.Environment.GetId()); err != nil {
		return err
	}

	r.WithKubeConfigPath(kubeConfigPath)

	return nil
}
