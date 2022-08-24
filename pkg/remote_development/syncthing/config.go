package syncthing

import (
	"encoding/xml"
	"io"
	"os"
)

type Configuration struct {
	XMLName xml.Name `xml:"configuration"`
	Version int      `xml:"version,attr"`
	Folder  []Folder `xml:"folder"`
	Device  []Device `xml:"device"`
	GUI     GUI      `xml:"gui"`
}

type Folder struct {
	XMLName                 xml.Name       `xml:"folder"`
	Id                      string         `xml:"id,attr"`
	Label                   string         `xml:"label,attr"`
	Path                    string         `xml:"path,attr"`
	Type                    string         `xml:"type,attr"`
	RescanIntervalS         int            `xml:"rescanIntervalS,attr"`
	FsWatcherEnabled        string         `xml:"fsWatcherEnabled,attr"`
	FsWatcherDelayS         int            `xml:"fsWatcherDelayS,attr"`
	IgnorePerms             string         `xml:"ignorePerms,attr"`
	AutoNormalize           string         `xml:"autoNormalize,attr"`
	FilesystemType          string         `xml:"filesystemType"`
	Device                  []FolderDevice `xml:"device"`
	MinDiskFree             MinDiskFree    `xml:"minDiskFree"`
	Versioning              Versioning     `xml:"versioning"`
	Copiers                 int            `xml:"copiers"`
	PullerMaxPendingKiB     string         `xml:"pullerMaxPendingKiB"`
	Hashers                 int            `xml:"hashers"`
	Order                   string         `xml:"order"`
	IgnoreDelete            bool           `xml:"ignoreDelete"`
	ScanProgressIntervalS   int            `xml:"scanProgressIntervalS"`
	PullerPauseS            int            `xml:"pullerPauseS"`
	MaxConflicts            int            `xml:"maxConflicts"`
	DisableSparseFiles      bool           `xml:"disableSparseFiles"`
	DisableTempIndexes      bool           `xml:"disableTempIndexes"`
	Paused                  bool           `xml:"paused"`
	WeakHashThresholdPct    int            `xml:"WeakHashThresholdPct"`
	MarkerName              string         `xml:"markerName"`
	CopyOwnershipFromParent bool           `xml:"copyOwnershipFromParent"`
	ModTimeWindowS          int            `xml:"modTimeWindowS"`
	MaxConcurrentWrites     int            `xml:"maxConcurrentWrites"`
	DisableFsync            bool           `xml:"disableFsync"`
	BlockPullOrder          string         `xml:"blockPullOrder"`
	CopyRangeMethod         string         `xml:"copyRangeMethod"`
	CaseSensitiveFS         bool           `xml:"caseSensitiveFS"`
	JunctionsAsDirs         bool           `xml:"junctionsAsDirs"`
}

type FolderDevice struct {
	XMLName            xml.Name `xml:"device"`
	Id                 string   `xml:"id,attr"`
	IntroducedBy       string   `xml:"introducedBy,attr"`
	EncryptionPassword string   `xml:"encryptionPassword"`
}

type MinDiskFree struct {
	Unit  string `xml:"unit,attr"`
	Value int    `xml:",chardata"`
}

type Versioning struct {
	XMLName          xml.Name `xml:"versioning"`
	CleanupIntervalS int      `xml:"cleanupIntervalS"`
	FsPath           string   `xml:"fsPath"`
	FsType           string   `xml:"fsType"`
}

type Device struct {
	XMLName                  xml.Name `xml:"device"`
	Id                       string   `xml:"id,attr"`
	Name                     string   `xml:"name,attr"`
	Compression              string   `xml:"compression,attr"`
	Introducer               bool     `xml:"introducer,attr"`
	SkipIntroductionRemovals bool     `xml:"skipIntroductionRemovals,attr"`
	IntroducedBy             string   `xml:"introducedBy,attr"`
	Address                  string   `xml:"address"`
	Paused                   bool     `xml:"paused"`
	AutoAcceptFolders        bool     `xml:"autoAcceptFolders"`
	MaxSendKbps              int      `xml:"maxSendKbps"`
	MaxRecvKbps              int      `xml:"maxRecvKbps"`
	MaxRequestKiB            int      `xml:"maxRequestKiB"`
	Untrusted                bool     `xml:"untrusted"`
	RemoteGUIPort            int      `xml:"remoteGUIPort"`
}

type GUI struct {
	XMLName   xml.Name `xml:"gui"`
	Enabled   bool     `xml:"enabled,attr"`
	TLS       bool     `xml:"tls,attr"`
	Debugging bool     `xml:"debugging,attr"`
	Address   string   `xml:"address"`
	ApiKey    string   `xml:"apikey"`
	Theme     string   `xml:"theme"`
	User      string   `xml:"user,omitempty"`
	Password  string   `xml:"password,omitempty"`
}

func LoadConfigXml(filePath string) (Configuration, error) {
	var config Configuration

	file, err := os.Open(filePath)
	if err != nil {
		return config, err
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return config, err
	}

	err = xml.Unmarshal(byteValue, &config)

	return config, err
}

func DumpConfigXml(config Configuration, filePath string) error {
	data, err := xml.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}
