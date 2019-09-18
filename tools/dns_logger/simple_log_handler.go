package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

type simpleLogHandler struct {
}

// constructor for symmetry with newMultiplexingLogHandler
func newSimpleLogHandler() *simpleLogHandler {
	return &simpleLogHandler{}
}

func (h *simpleLogHandler) Handle(domain string) error {
	parts := strings.Split(domain, ".")
	if len(parts) < 3 {
		return fmt.Errorf("not enough subdomains: %s", domain)
	}

	message, err := hex.DecodeString(parts[0])
	if err != nil {
		return fmt.Errorf("non-hex message")
	}

	os.Stdout.Write(message)

	return nil
}
