package config

type Context struct {
	Organization     string `json:"organization,omitempty" yaml:"organization,omitempty"`
	Project          string `json:"project,omitempty" yaml:"project,omitempty"`
	Environment      string `json:"environment,omitempty" yaml:"environment,omitempty"`
	ServiceComponent string `json:"serviceComponent,omitempty" yaml:"serviceComponent,omitempty"`
}

type Profile struct {
	Name string `json:"-" yaml:"-"`

	Host  string `json:"host,omitempty" yaml:"host,omitempty"`
	Token string `json:"token,omitempty" yaml:"token,omitempty"`

	Context Context `json:"context,omitempty" yaml:"context,omitempty"`
}

type NamedProfiles map[string]Profile
