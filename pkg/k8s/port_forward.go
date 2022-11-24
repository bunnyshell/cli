package k8s

type PortForward struct {
	Interface string

	RemotePort int
	LocalPort  int

	StopChannel  chan struct{}
	ReadyChannel chan struct{}
}

func NewPortForward(iface string, localPort int, remotePort int) *PortForward {
	return &PortForward{
		Interface: iface,

		RemotePort: remotePort,
		LocalPort:  localPort,

		StopChannel:  make(chan struct{}),
		ReadyChannel: make(chan struct{}, 1),
	}
}

func (p *PortForward) Close() {
	if p.StopChannel != nil {
		close(p.StopChannel)
	}
}
