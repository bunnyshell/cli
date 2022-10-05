package remote_development

import (
	"fmt"

	"bunnyshell.com/dev/pkg/remote"
	bunnysdk "bunnyshell.com/sdk"
)

var (
	ErrNoOrganizationSelected = fmt.Errorf("you need to select an organization first")
)

type RemoteDevelopment struct {
	remoteDev *remote.RemoteDevelopment

	organization *bunnysdk.OrganizationItem
	project      *bunnysdk.ProjectItem
	environment  *bunnysdk.EnvironmentItem
	component    *bunnysdk.ComponentItem

	environmentWorkspaceDir string

	kubeConfigPath string

	localSyncPath string
}

func NewRemoteDevelopment() *RemoteDevelopment {
	return &RemoteDevelopment{
		remoteDev: remote.NewRemoteDevelopment(),

		organization: nil,
		project:      nil,
		environment:  nil,
		component:    nil,
	}
}

func (r *RemoteDevelopment) WithOrganization(organization *bunnysdk.OrganizationItem) *RemoteDevelopment {
	r.organization = organization
	return r
}

func (r *RemoteDevelopment) WithProject(project *bunnysdk.ProjectItem) *RemoteDevelopment {
	if r.organization != nil && r.organization.GetId() != project.GetOrganization() {
		panic(fmt.Errorf(
			"project \"%s\" is not part of organization \"%s\"",
			project.GetName(),
			r.organization.GetName(),
		))
	}

	r.project = project
	return r
}

func (r *RemoteDevelopment) WithEnvironment(environment *bunnysdk.EnvironmentItem) *RemoteDevelopment {
	if r.project != nil && r.project.GetId() != environment.GetProject() {
		panic(fmt.Errorf(
			"environment \"%s\" is not part of project \"%s\"",
			environment.GetName(),
			r.project.GetName(),
		))
	}

	r.environment = environment
	return r
}

func (r *RemoteDevelopment) WithComponent(component *bunnysdk.ComponentItem) *RemoteDevelopment {
	if r.environment == nil {
		panic(fmt.Errorf("you have to select an environment before selecting a component"))
	}

	if r.environment.GetId() != component.GetEnvironment() {
		panic(fmt.Errorf(
			"component \"%s\" is not part of environment \"%s\"",
			component.GetName(),
			r.environment.GetName(),
		))
	}

	if component.GetSyncPath() == "" {
		panic(fmt.Errorf("component has no syncPath defined"))
	}

	r.component = component
	return r
}

func (r *RemoteDevelopment) WithEnvironmentWorkspaceDir(environmentWorkspaceDir string) *RemoteDevelopment {
	r.environmentWorkspaceDir = environmentWorkspaceDir
	return r
}

func (r *RemoteDevelopment) WithKubeConfigPath(kubeConfigPath string) *RemoteDevelopment {
	r.kubeConfigPath = kubeConfigPath
	return r
}

func (r *RemoteDevelopment) WithLocalSyncPath(localSyncPath string) *RemoteDevelopment {
	r.localSyncPath = localSyncPath
	return r
}
