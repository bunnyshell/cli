package ssh

import (
	"io"
	"log"
	"net"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"
)

type SSHTunnel struct {
	Local  *Endpoint
	Server *Endpoint
	Remote *Endpoint

	Config *ssh.ClientConfig

	Log *log.Logger

	ReadyChannel chan bool
	StopChannel  chan bool

	listener net.Listener
}

func (tunnel *SSHTunnel) logf(fmt string, args ...interface{}) {
	if tunnel.Log != nil {
		tunnel.Log.Printf(fmt, args...)
	}
}

func (tunnel *SSHTunnel) Start() error {
	if err := tunnel.listen(); err != nil {
		return err
	}

	tunnel.Local.Port = tunnel.listener.Addr().(*net.TCPAddr).Port
	if tunnel.ReadyChannel != nil {
		close(tunnel.ReadyChannel)
	}

	<-tunnel.StopChannel

	return nil
}

func (tunnel *SSHTunnel) listen() error {
	listener, err := net.Listen("tcp", tunnel.Local.String())
	if err != nil {
		return err
	}
	tunnel.listener = listener

	go tunnel.waitForConnection()

	return nil
}

func (tunnel *SSHTunnel) waitForConnection() {
	for {
		conn, err := tunnel.listener.Accept()
		if err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				tunnel.logf("error on listener.Accept: %s", err)
			}
			return
		}

		go tunnel.handleConnection(conn)
	}
}

func (tunnel *SSHTunnel) handleConnection(localConn net.Conn) {
	defer localConn.Close()

	serverConn, err := ssh.Dial("tcp", tunnel.Server.String(), tunnel.Config)
	if err != nil {
		tunnel.logf("server dial error: %s", err)
		return
	}
	defer serverConn.Close()

	remoteConn, err := serverConn.Dial("tcp", tunnel.Remote.String())
	if err != nil {
		tunnel.logf("remote dial error: %s", err)
		return
	}
	defer remoteConn.Close()

	var wg sync.WaitGroup
	copyConn := func(writer, reader net.Conn) {
		defer wg.Done()
		_, err := io.Copy(writer, reader)
		if err != nil {
			tunnel.logf("io.Copy error: %s", err)
		}
	}

	wg.Add(1)
	go copyConn(localConn, remoteConn)

	wg.Add(1)
	go copyConn(remoteConn, localConn)

	wg.Wait()
}

func (tunnel *SSHTunnel) Stop() {
	if tunnel.listener != nil {
		tunnel.listener.Close()
	}

	if tunnel.StopChannel != nil {
		close(tunnel.StopChannel)
	}
}

func NewSSHTunnel(serverHost string, serverPort int, auth ssh.AuthMethod, remoteHost string, remotePort int) *SSHTunnel {
	localEndpoint := NewEndpoint("127.0.0.1", 0)
	server := NewEndpoint(serverHost, serverPort)

	return &SSHTunnel{
		Config: &ssh.ClientConfig{
			User: server.User,
			Auth: []ssh.AuthMethod{auth},
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
		},
		// @todo replace with a logger, once we have one in cli
		Log:          nil,
		Local:        localEndpoint,
		Server:       server,
		Remote:       NewEndpoint(remoteHost, remotePort),
		ReadyChannel: make(chan bool),
		StopChannel:  make(chan bool),
	}
}
