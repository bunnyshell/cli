package ssh

import "fmt"

type Endpoint struct {
	Host string
	Port int
	User string
}

func NewEndpoint(host string, port int) *Endpoint {
	return &Endpoint{
		Host: host,
		Port: port,
	}
}

func (e *Endpoint) WithUser(user string) *Endpoint {
	e.User = user
	return e
}

func (e *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", e.Host, e.Port)
}
