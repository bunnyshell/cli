package action

import (
	"os"

	"bunnyshell.com/cli/pkg/interactive"
	mutagenConfig "bunnyshell.com/dev/pkg/mutagen/config"
	"bunnyshell.com/dev/pkg/util"
)

func (up *Up) ensureSyncPaths(parameters *UpParameters) error {
	if parameters.SyncMode == mutagenConfig.None {
		return up.ensurePersistentWorkdir(parameters)
	}

	if err := up.ensureLocalSyncPath(parameters); err != nil {
		return err
	}

	if err := up.ensureRemoteSyncPath(parameters); err != nil {
		return err
	}

	return nil
}

func (up *Up) ensureLocalSyncPath(parameters *UpParameters) error {
	if parameters.LocalSyncPath != "" {
		return nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	question := interactive.NewInput("Local Sync Path")
	question.Default = cwd
	question.Help = "Local path is the folder on your machine that will be synced into the container"
	question.SetValidate(util.IsDirectoryValidator)

	syncPath, err := question.AskString()
	if err != nil {
		return err
	}

	parameters.LocalSyncPath = syncPath

	return nil
}

func (up *Up) ensureRemoteSyncPath(parameters *UpParameters) error {
	if parameters.RemoteSyncPath != "" {
		return nil
	}

	question := interactive.NewInput("Remote Sync Path")
	question.Help = "Remote path is the folder within the container where the application is loaded from\n" +
		"This is where the local files will be synced to\n" +
		"This folder will be mounted as a persistent volume to persist your changes across multiple development sessions."

	syncPath, err := question.AskString()
	if err != nil {
		return err
	}

	parameters.RemoteSyncPath = syncPath

	return nil
}

func (up *Up) ensurePersistentWorkdir(parameters *UpParameters) error {
	if parameters.RemoteSyncPath != "" {
		return nil
	}

	question := interactive.NewInput("Persistent Workdir")
	question.Help = "Persistent workdir is the folder within the container where the application is loaded from\n" +
		"This folder will be mounted as a persistent volume to persist your changes across multiple development sessions."

	syncPath, err := question.AskString()
	if err != nil {
		return err
	}

	parameters.RemoteSyncPath = syncPath

	return nil
}
