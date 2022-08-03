package remote_development

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"bunnyshell.com/cli/pkg/remote_development/syncthing"
	"bunnyshell.com/cli/pkg/util"
)

const (
	SyncthingVersion             = "v1.20.4"
	SyncthingLocalConfigDirname  = "syncthing-local"
	SyncthingRemoteConfigDirname = "syncthing-remote"

	syncthingBinFilename      = "syncthing"
	syncthingDownloadFilename = "syncthing-%s-%s-%s.%s"
	syncthingDownloadUrl      = "https://github.com/syncthing/syncthing/releases/download/%s/%s"
)

func (r *RemoteDevelopment) PrepareSyncthing() error {
	spinner := util.MakeSpinner(" Prepare Syncthing")
	spinner.Start()
	defer spinner.Stop()

	syncthingBinPath, err := ensureSyncthingBin()
	if err != nil {
		return err
	}

	localConfigDir := r.getLocalSyncthingConfigDir()
	generateLocalConfigCommand := exec.Command(syncthingBinPath, "generate", "--home", localConfigDir)
	output, err := generateLocalConfigCommand.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", output)
	}

	remoteConfigDir := r.getRemoteSyncthingConfigDir()
	generateRemoteConfigCommand := exec.Command(syncthingBinPath, "generate", "--home", remoteConfigDir)
	output, err = generateRemoteConfigCommand.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", output)
	}

	return patchSyncthingConfigs(localConfigDir, remoteConfigDir, r.LocalSyncPath, r.RemoteSyncPath)
}

func (r *RemoteDevelopment) UpdateLocalSyncthingConfig() error {
	localConfigDir := r.getLocalSyncthingConfigDir()
	localConfigPath := localConfigDir + "/config.xml"
	localConfig, err := syncthing.LoadConfigXml(localConfigPath)
	if err != nil {
		return err
	}

	for idx := range localConfig.Device {
		device := &localConfig.Device[idx]
		if device.Name == "remote" {
			device.Address = fmt.Sprintf("tcp://%s", r.SyncthingSSHTunnel.Local.String())
		}
	}

	return syncthing.DumpConfigXml(localConfig, localConfigPath)
}

// @todo add restart/stop management
func (r *RemoteDevelopment) StartLocalSyncthing() error {
	syncthingBinPath, err := getSyncthingBinPath()
	if err != nil {
		return err
	}

	localConfigDir := r.getLocalSyncthingConfigDir()
	logFilePath := filepath.Join(filepath.Dir(syncthingBinPath), "syncthing.log")
	logMaxSizeBytes := 10 * 1024 * 1024
	syncthingArgs := []string{"serve", "--home", localConfigDir, "--no-browser", "--logfile", logFilePath, "--log-max-old-files", "0", "--log-max-size", strconv.Itoa(logMaxSizeBytes)}
	syncthingCmd := exec.Command(syncthingBinPath, syncthingArgs...)
	if err := syncthingCmd.Start(); err != nil {
		return err
	}

	go syncthingCmd.Wait()
	return nil
}

func (r *RemoteDevelopment) getLocalSyncthingConfigDir() string {
	return fmt.Sprintf("%s/%s", r.ComponentFolderPath, SyncthingLocalConfigDirname)
}

func (r *RemoteDevelopment) getRemoteSyncthingConfigDir() string {
	return fmt.Sprintf("%s/%s", r.ComponentFolderPath, SyncthingRemoteConfigDirname)
}

func patchSyncthingConfigs(localConfigDir, remoteConfigDir, localSyncFolder, remoteSyncFolder string) error {
	localConfigPath := localConfigDir + "/config.xml"
	localConfig, err := syncthing.LoadConfigXml(localConfigPath)
	if err != nil {
		return err
	}

	remoteConfigPath := remoteConfigDir + "/config.xml"
	remoteConfig, err := syncthing.LoadConfigXml(remoteConfigPath)
	if err != nil {
		return err
	}

	// update configs
	localConfig.GUI.Enabled = false
	remoteConfig.GUI.Enabled = false

	folderLocal := &localConfig.Folder[0]
	folderLocal.Id = "sync-folder"
	folderLocal.Label = "Sync Folder"
	folderLocal.Path = localSyncFolder
	folderLocal.Type = "sendonly"
	folderLocal.RescanIntervalS = 300
	folderLocal.FsWatcherDelayS = 1
	folderLocal.MarkerName = "."
	folderLocal.MaxConflicts = 0

	folderRemote := &remoteConfig.Folder[0]
	folderRemote.Id = "sync-folder"
	folderRemote.Label = "Sync Folder"
	folderRemote.Path = remoteSyncFolder
	folderRemote.RescanIntervalS = 300
	folderRemote.FsWatcherDelayS = 1
	folderRemote.MarkerName = "."
	folderRemote.MaxConflicts = 0

	deviceLocal := &localConfig.Device[0]
	deviceLocal.Name = "local"

	deviceRemote := &remoteConfig.Device[0]
	deviceRemote.Name = "remote"

	// sync devices
	if len(localConfig.Device) < 2 {
		localConfig.Device = append(localConfig.Device, *deviceRemote)
		remoteConfig.Device = append(remoteConfig.Device, *deviceLocal)
		folderLocal.Device = append(folderLocal.Device, folderRemote.Device[0])
		folderRemote.Device = append(folderRemote.Device, folderLocal.Device[0])
	}

	err = syncthing.DumpConfigXml(localConfig, localConfigPath)
	if err != nil {
		return err
	}
	err = syncthing.DumpConfigXml(remoteConfig, remoteConfigPath)
	if err != nil {
		return err
	}

	return nil
}

func getSyncthingBinPath() (string, error) {
	workspaceDir, err := util.GetWorkspaceDir()
	if err != nil {
		return "", err
	}

	return workspaceDir + "/" + syncthingBinFilename, nil
}

func ensureSyncthingBin() (string, error) {
	syncthingBinPath, err := getSyncthingBinPath()
	if err != nil {
		return "", err
	}

	stats, err := os.Stat(syncthingBinPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return syncthingBinPath, err
	}
	if err == nil && stats.Size() > 0 && !stats.IsDir() {
		return syncthingBinPath, nil
	}

	downloadFilename := fmt.Sprintf(syncthingDownloadFilename, getPlatform(), runtime.GOARCH, SyncthingVersion, getExtension())
	syncthingArchivePath := filepath.Dir(syncthingBinPath) + "/" + downloadFilename
	downloadUrl := fmt.Sprintf(syncthingDownloadUrl, SyncthingVersion, downloadFilename)

	err = downloadSyncthingArchive(downloadUrl, syncthingArchivePath)
	if err != nil {
		return syncthingBinPath, err
	}

	err = extractSyncthingBin(syncthingArchivePath, syncthingBinPath)
	if err != nil {
		return syncthingBinPath, err
	}

	return syncthingBinPath, removeSyncthingArchive(syncthingArchivePath)
}

func removeSyncthingArchive(filePath string) error {
	return os.Remove(filePath)
}

func downloadSyncthingArchive(source, destination string) error {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	out, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := client.Get(source)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func extractSyncthingBin(source, destination string) error {
	if runtime.GOOS == "darwin" || runtime.GOOS == "windows" {
		return extractSyncthingBinZip(source, destination)
	}

	return extractSyncthingBinTarGz(source, destination)
}

func extractSyncthingBinZip(source, destination string) error {
	reader, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, f := range reader.File {
		if strings.Split(f.Name, "/")[1] == "syncthing" {
			destinationFile, err := os.OpenFile(destination, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer destinationFile.Close()

			zippedFile, err := f.Open()
			if err != nil {
				return err
			}
			defer zippedFile.Close()

			if _, err := io.Copy(destinationFile, zippedFile); err != nil {
				return err
			}
			return nil
		}
	}

	return nil
}

func extractSyncthingBinTarGz(source, destination string) error {
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}

	gzipReader, err := gzip.NewReader(sourceFile)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(gzipReader)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		if strings.Split(header.Name, "/")[1] == "syncthing" {
			destinationFile, err := os.OpenFile(destination, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, header.FileInfo().Mode())
			if err != nil {
				return err
			}
			defer destinationFile.Close()

			if _, err := io.Copy(destinationFile, tarReader); err != nil {
				return err
			}
			return nil
		}
	}

	return nil
}

func getExtension() string {
	if runtime.GOOS == "darwin" || runtime.GOOS == "windows" {
		return "zip"
	}

	return "tar.gz"
}

func getPlatform() string {
	if runtime.GOOS == "darwin" {
		return "macos"
	}

	return runtime.GOOS
}
