package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"time"
)

type multiplexingLogHandler struct {
	msgs                 map[string]chan []byte
	dnsMessageIDLength   int
	interFragmentTimeout time.Duration
	*sync.Mutex
}

func newMultiplexingLogHandler(idLength int, interFragmentTimeout time.Duration) *multiplexingLogHandler {
	return &multiplexingLogHandler{
		msgs:                 make(map[string]chan []byte),
		dnsMessageIDLength:   idLength,
		interFragmentTimeout: interFragmentTimeout,
		Mutex:                new(sync.Mutex),
	}
}

func (h *multiplexingLogHandler) Handle(domain string) error {
	fmt.Println(domain)
	id, hexFragment, err := h.parseFragment(domain)
	if err != nil {
		return err
	}

	if hexFragment == "close" {
		_, ok := h.msgs[id]
		if ok {
			close(h.msgs[id])
			h.Lock()
			delete(h.msgs, id)
			h.Unlock()
		}
		return nil
	}

	frag, err := hex.DecodeString(hexFragment)
	if err != nil {
		return fmt.Errorf("non-hex fragment")
	}

	h.Lock()
	_, ok := h.msgs[id]
	h.Unlock()
	if !ok {
		c := make(chan []byte)
		h.Lock()
		h.msgs[id] = c
		h.Unlock()
		go h.messageWorker(c)
	}

	h.Lock()
	h.msgs[id] <- frag
	h.Unlock()

	return nil
}

func (h *multiplexingLogHandler) messageWorker(input chan []byte) {
	buf := bytes.Buffer{}

loop:
	for {
		select {
		case fragment, ok := <-input:
			if !ok {
				// message is over
				break loop
			}
			buf.Write(fragment)
		case <-time.After(h.interFragmentTimeout):
			fmt.Printf("Message timeout: %s\n", buf.String())
			return
		}
	}

	fmt.Println(buf.String())

}

func (h *multiplexingLogHandler) parseFragment(domain string) (id string, fragment string, err error) {
	parts := strings.Split(domain, ".")
	if len(parts) < 4 {
		return "", "", fmt.Errorf("not enough subdomains: %s", domain)
	}
	if len(parts[1]) != h.dnsMessageIDLength {
		return "", "", fmt.Errorf("id field is not of length %d", h.dnsMessageIDLength)
	}
	return parts[1], parts[0], nil
}
