package gateway

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/hashicorp/yamux"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

// Opts hold the options for the gateway command
type Opts struct {
	connectBackAddress string
	gatewayAddress     string
}

var gatewayOpts Opts

// GatewayCmd contains the CLI interface for the gateway command.
var GatewayCmd = &cobra.Command{
	Use:   "gateway",
	Short: "A gateway for the reverse socks5 payload",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runGateway(gatewayOpts)
	},
}

func init() {
	flags := GatewayCmd.PersistentFlags()
	flags.StringVarP(&gatewayOpts.connectBackAddress, "connectback", "c", "",
		"Connect back `listen address` to which the payload will connect")
	flags.StringVarP(&gatewayOpts.gatewayAddress, "gateway", "g", "",
		"Gateway `listen address` which exposes the payloads socks5 proxy")

	_ = GatewayCmd.MarkPersistentFlagRequired("connectback")
	_ = GatewayCmd.MarkPersistentFlagRequired("gateway")
}

func runGateway(opts Opts) error {
	listener, err := net.Listen("tcp", opts.connectBackAddress)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	log.Printf("listening for payload connection on tcp/%s", opts.connectBackAddress)

	targetConn, err := listener.Accept()
	if err != nil {
		return fmt.Errorf("accept: %w", err)
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
		return fmt.Errorf("listen: %w", err)
	}

	log.Printf("start socks5 gateway on tcp/%s", gatewayAddress)

	for {
		conn, err := gateway.Accept()
		if err != nil {
			return fmt.Errorf("accept: %w", err)
		}

		go func() {
			err := handleGatewayConn(conn, sess)
			if err != nil {
				log.Printf("error while handling gateway connection: %v", err)
			}
		}()
	}
}

func handleGatewayConn(conn net.Conn, sess *yamux.Session) error {
	yamuxConn, err := sess.Open()
	if err != nil {
		return fmt.Errorf("open multiplexed connection: %w", err)
	}

	log.Print("multiplexed connection opened")

	eg, _ := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		_, err := io.Copy(yamuxConn, conn)
		return err
	})
	eg.Go(func() error {
		_, err := io.Copy(conn, yamuxConn)
		return err
	})

	return eg.Wait()
}
