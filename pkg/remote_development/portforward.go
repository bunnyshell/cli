package remote_development

import (
	"log"
	"os"

	bunnyshellSSH "bunnyshell.com/cli/pkg/ssh"
	"bunnyshell.com/cli/pkg/util"
)

func (r *RemoteDevelopment) EnsureRemoteSSHPortForward() error {
	spinner := util.MakeSpinner(" Start Remote SSH Port Forward")
	spinner.Start()
	defer spinner.Stop()

	forwarder, err := r.KubernetesClient.PortForwardRemoteSSH(r.ComponentName)
	if err != nil {
		return err
	}
	r.RemoteSSHForwarder = forwarder

	return nil
}

func (r *RemoteDevelopment) EnsureSyncthingPortForward() error {
	tunnel := bunnyshellSSH.NewSSHTunnel(
		r.KubernetesClient.SSHPortForwardOptions.Interface,
		r.KubernetesClient.SSHPortForwardOptions.LocalPort,
		bunnyshellSSH.PrivateKeyFile(r.SSHPrivateKeyPath),
		SyncthingRemoteInterface,
		SyncthingRemotePort,
	)
	// todo replace with a general logging solution
	tunnel.Log = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)

	errChan := make(chan error, 1)
	go func() {
		errChan <- tunnel.Start()
		close(errChan)
	}()

	select {
	case <-tunnel.ReadyChannel:
		r.SyncthingSSHTunnel = tunnel
	case err := <-errChan:
		return err
	}

	return nil
}
