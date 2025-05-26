package workspace

import (
	"os"
	"path/filepath"

	"bunnyshell.com/cli/pkg/api/environment"
	"bunnyshell.com/cli/pkg/util"
)

type Workspace struct {
	environment string

	rootDir string

	kubeConfigFile string
}

func NewWorkspace(environment string) *Workspace {
	return &Workspace{
		environment: environment,
	}
}

func (workspace *Workspace) GetEnvironmentDir() (string, error) {
	if workspace.rootDir == "" {
		dir, err := workspace.makeWorkspace()
		if err != nil {
			return "", err
		}

		workspace.rootDir = dir
	}

	return workspace.rootDir, nil
}

func (workspace *Workspace) GetKubeConfigFile() (string, error) {
	if workspace.kubeConfigFile == "" {
		envWorkspace, err := workspace.GetEnvironmentDir()
		if err != nil {
			return "", err
		}

		workspace.kubeConfigFile = filepath.Join(envWorkspace, kubeConfigFilename)
	}

	return workspace.kubeConfigFile, nil
}

func (workspace *Workspace) DownloadKubeConfig(overrideClusterServer string) (string, error) {
	kubeConfigFile, err := workspace.GetKubeConfigFile()
	if err != nil {
		return "", err
	}

	kubeConfig, err := environment.KubeConfig(
		environment.NewKubeConfigOptions(workspace.environment, overrideClusterServer),
	)
	if err != nil {
		return "", err
	}

	if err = os.WriteFile(kubeConfigFile, kubeConfig.Bytes, kubeConfigPerm); err != nil {
		return "", err
	}

	return kubeConfigFile, nil
}

func (workspace *Workspace) makeWorkspace() (string, error) {
	workspaceRoot, err := util.GetWorkspaceDir()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(workspaceRoot, workspace.environment)

	if err = os.MkdirAll(dir, workspacePerm); err != nil {
		return "", err
	}

	return dir, nil
}
