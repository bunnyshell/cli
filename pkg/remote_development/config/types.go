package config

type SyncPath struct {
	RemotePath string `yaml:"remotePath"`
	LocalPath  string `yaml:"localPath"`
}

type Environ map[string]string

type (
	PortForward     string
	PortForwardList []PortForward
)

type Resource struct {
	CPU    string `yaml:"cpu,omitempty"`
	Memory string `yaml:"memory,omitempty"`
}

type ResourceList struct {
	Limits   Resource `yaml:"limits,omitempty"`
	Requests Resource `yaml:"requests,omitempty"`
}

type Profile struct {
	Name string `yaml:"-"`

	Command []string `yaml:"command,omitempty"`

	SyncPaths []SyncPath `yaml:"syncPaths,omitempty"`

	PortForwards PortForwardList `yaml:"portForwards,omitempty"`

	Environ Environ `yaml:"environment,omitempty"`

	Resources ResourceList `yaml:"resources,omitempty"`
}

type NamedProfiles map[string]Profile
