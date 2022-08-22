package remote_development

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/ssh"

	bunnysdk "bunnyshell.com/sdk"
	"k8s.io/client-go/tools/portforward"
)

type RemoteDevelopment struct {
	SSHPrivateKeyPath  string
	SSHPublicKeyPath   string
	RemoteSSHForwarder *portforward.PortForwarder
	SyncthingSSHTunnel *ssh.SSHTunnel

	SyncthingProcess *os.Process

	OrganizationId string
	ProjectId      string
	EnvironmentId  string

	LocalSyncPath  string
	RemoteSyncPath string

	ComponentId         string
	ComponentName       string
	ComponentFolderPath string

	KubeConfigPath   string
	KubernetesClient *KubernetesClient
	ContainerName    string

	StopChannel chan bool
}

func NewRemoteDevelopment() *RemoteDevelopment {
	return &RemoteDevelopment{
		StopChannel: make(chan bool),
	}
}

func (r *RemoteDevelopment) WithSSH(sshPrivateKeyPath, sshPublicKeyPath string) *RemoteDevelopment {
	r.SSHPrivateKeyPath = sshPrivateKeyPath
	r.SSHPublicKeyPath = sshPublicKeyPath
	return r
}

func (r *RemoteDevelopment) WithKubernetesClient(kubeConfigPath string) *RemoteDevelopment {
	r.KubeConfigPath = kubeConfigPath
	kubernetesClient, err := NewKubernetesClient(kubeConfigPath)
	if err != nil {
		panic(err)
	}

	r.KubernetesClient = kubernetesClient

	return r
}

type ComponentItem interface {
	GetId() string
	GetName() string
	GetSyncPath() string
}

func (r *RemoteDevelopment) WithComponent(component ComponentItem) *RemoteDevelopment {
	r.ComponentId = component.GetId()
	r.ComponentName = component.GetName()
	r.RemoteSyncPath = component.GetSyncPath()

	return r
}

func (r *RemoteDevelopment) StartRemoteDevelopment() error {
	if _, _, err := r.componentRemoteDevelopmentUp(); err != nil {
		return fmt.Errorf("remote development up failed: %s", err.Error())
	}

	// wait for pod to be ready
	if err := r.KubernetesClient.WatchRemoteDevPods(r.ComponentName, r.ContainerName); err != nil {
		return err
	}

	if err := r.EnsureRemoteSSHPortForward(); err != nil {
		return err
	}

	if err := r.EnsureSSHConfigEntry(); err != nil {
		return err
	}

	if err := r.EnsureSyncthingPortForward(); err != nil {
		return err
	}

	if err := r.UpdateLocalSyncthingConfig(); err != nil {
		return err
	}

	if err := r.StartLocalSyncthing(); err != nil {
		return err
	}

	return r.StartSSHTerminal()
}

func (r *RemoteDevelopment) StopRemoteDevelopment() error {
	if _, _, err := r.componentRemoteDevelopmentDown(); err != nil {
		return fmt.Errorf("remote development down failed: %s", err.Error())
	}

	return nil
}

func (r *RemoteDevelopment) Close() {
	// close syncthing ssh tunnel
	if r.SyncthingSSHTunnel != nil {
		r.SyncthingSSHTunnel.Stop()
	}

	// close k8s remote ssh portforwarding
	if r.RemoteSSHForwarder != nil {
		r.RemoteSSHForwarder.Close()
	}

	// close cli command
	if r.StopChannel != nil {
		close(r.StopChannel)
	}

	// kill syncthing process
	if r.SyncthingProcess != nil {
		r.SyncthingProcess.Kill()
	}
}

func (r *RemoteDevelopment) CloseOnSignal(signals chan os.Signal) {
	<-signals
	r.Close()
}

func (r *RemoteDevelopment) Wait() {
	<-r.StopChannel
}

func (r *RemoteDevelopment) componentRemoteDevelopmentUp() (*bunnysdk.ComponentItem, *http.Response, error) {
	ctx, cancel := lib.GetContext()
	defer cancel()

	request := lib.GetAPI().ComponentApi.ComponentRemoteDevelopmentUp(ctx, r.ComponentId).ComponentRemoteDevelopmentUp(
		bunnysdk.ComponentRemoteDevelopmentUp{
			Container: *bunnysdk.NewNullableString(&r.ContainerName),
		},
	)
	return request.Execute()
}

func (r *RemoteDevelopment) componentRemoteDevelopmentDown() (*bunnysdk.ComponentItem, *http.Response, error) {
	ctx := context.WithValue(context.Background(), bunnysdk.ContextAPIKeys, map[string]bunnysdk.APIKey{
		"ApiKeyAuth": {
			Key: lib.CLIContext.Profile.Token,
		},
	})

	ctx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	request := lib.GetAPI().ComponentApi.ComponentRemoteDevelopmentDown(ctx, r.ComponentId)
	return request.Execute()
}
