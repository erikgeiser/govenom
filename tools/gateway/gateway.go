package gateway

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/hashicorp/yamux"
	"golang.org/x/sync/errgroup"
)

func runGatwayForPayloadConnection(opts Opts) error {
	listener, err := net.Listen("tcp", opts.connectBackAddress)
	if err != nil {
		return fmt.Errorf("listen for payload connection: %w", err)
	}

	log.Printf("waiting for payload connection on %s/tcp", opts.connectBackAddress)

	targetConn, err := listener.Accept()
	if err != nil {
		return fmt.Errorf("accept payload connection: %w", err)
	}

	err = listener.Close()
	if err != nil {
		return fmt.Errorf("closing payload listener: %w", err)
	}

	log.Printf("received connection from %s", targetConn.RemoteAddr())

	client, err := yamux.Client(targetConn, nil)
	if err != nil {
		return fmt.Errorf("connection multiplexer setup: %w", err)
	}

	return startGateway(opts.gatewayAddress, client)
}

func startGateway(gatewayAddress string, sess *yamux.Session) error {
	gateway, err := net.Listen("tcp", gatewayAddress)
	if err != nil {
		return fmt.Errorf("listen for gateway connection: %w", err)
	}

	var closedBecausePayloadDisconnected bool

	go func() {
		<-sess.CloseChan()

		log.Println("payload disconnected: stopping gateway")

		closedBecausePayloadDisconnected = true

		err := gateway.Close()
		if err != nil {
			log.Printf("gateway close: %v", err)
		}
	}()

	log.Printf("start socks5 gateway on %s/tcp", gatewayAddress)

	for {
		conn, err := gateway.Accept()
		if err != nil {
			if closedBecausePayloadDisconnected {
				return nil
			}

			return fmt.Errorf("accept gateway connection: %w", err)
		}

		go func() {
			err := handleGatewayConn(conn, sess)
			if err != nil {
				log.Printf("handling gateway connection: %v", err)
			}
		}()
	}
}

func handleGatewayConn(conn net.Conn, sess *yamux.Session) error {
	yamuxConn, err := sess.Open()
	if err != nil {
		return fmt.Errorf("open multiplexed connection: %w", err)
	}

	var eg errgroup.Group

	eg.Go(func() error {
		_, err := io.Copy(yamuxConn, conn)
		if err != nil {
			return fmt.Errorf("gateway->payload: %w", err)
		}

		return nil
	})
	eg.Go(func() error {
		_, err := io.Copy(conn, yamuxConn)
		if err != nil {
			return fmt.Errorf("payload->gateway: %w", err)
		}

		return nil
	})

	return eg.Wait()
}
