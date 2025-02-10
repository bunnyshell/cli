package debug_component

import (
	"os"
	"path/filepath"

	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/util"
)

const (
	KubeConfigFilename = "kube-config.yaml"
)

// @deprecated Use bunnyshell.com/cli/pkg/remote_development/workspace
func (d *DebugComponent) ensureEnvironmentWorkspaceDir() error {
	workspace, err := util.GetWorkspaceDir()
	if err != nil {
		return err
	}

	d.WithEnvironmentWorkspaceDir(filepath.Join(workspace, d.environmentResource.Environment.GetId()))

	return os.MkdirAll(d.environmentWorkspaceDir, 0755)
}

// @deprecated Use bunnyshell.com/cli/pkg/remote_development/workspace
func (d *DebugComponent) ensureEnvironmentKubeConfig() error {
	kubeConfigPath := filepath.Join(d.environmentWorkspaceDir, KubeConfigFilename)
	if err := lib.DownloadEnvironmentKubeConfig(kubeConfigPath, d.environmentResource.Environment.GetId()); err != nil {
		return err
	}

	d.WithKubeConfigPath(kubeConfigPath)

	return nil
}
