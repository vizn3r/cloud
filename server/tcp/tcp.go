package tcp

import "net"

type TCP struct {
	Listener net.Listener
	Host     string
}

func NewTCP(host string) *TCP {
	return &TCP{nil, host}
}

func (tcp *TCP) Start() error {
	var err error
	tcp.Listener, err = net.Listen("tcp", tcp.Host)

	return err
}
