package ssh

import (
	"fmt"
	"net"
	"os"
	"runtime"

	"github.com/shiena/ansicolor"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

type SSHTerminal struct {
	Server *Endpoint
	Config *ssh.ClientConfig

	ReadyChannel chan bool
}

func NewSSHTerminal(host string, port int, auth ssh.AuthMethod) *SSHTerminal {
	server := NewEndpoint(host, port)

	return &SSHTerminal{
		Config: &ssh.ClientConfig{
			User: server.User,
			Auth: []ssh.AuthMethod{auth},
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
		},
		Server: server,

		ReadyChannel: make(chan bool),
	}
}

func (sshTerminal *SSHTerminal) Start() error {
	serverConn, err := ssh.Dial("tcp", sshTerminal.Server.String(), sshTerminal.Config)
	if err != nil {
		return err
	}
	defer serverConn.Close()

	session, err := serverConn.NewSession()
	if err != nil {
		return err
	}

	session.Stdin = os.Stdin
	if runtime.GOOS == "windows" {
		session.Stdout = ansicolor.NewAnsiColorWriter(os.Stdout)
		session.Stderr = ansicolor.NewAnsiColorWriter(os.Stderr)
	} else {
		session.Stdout = os.Stdout
		session.Stderr = os.Stderr
	}

	termFd := int(os.Stdout.Fd())
	if !term.IsTerminal(termFd) {
		return fmt.Errorf("no terminal available")
	}

	if runtime.GOOS != "windows" {
		oldState, err := term.MakeRaw(termFd)
		if err != nil {
			return err
		}
		defer term.Restore(termFd, oldState)
	}

	w, h, err := term.GetSize(termFd)
	if err != nil {
		return err
	}

	terminalModes := ssh.TerminalModes{
		ssh.ECHO:  0,
		ssh.IGNCR: 1,
	}
	if err := session.RequestPty("xterm-256color", h, w, terminalModes); err != nil {
		return err
	}

	if err := session.Shell(); err != nil {
		return err
	}

	if sshTerminal.ReadyChannel != nil {
		close(sshTerminal.ReadyChannel)
	}

	session.Wait()
	return nil
}
