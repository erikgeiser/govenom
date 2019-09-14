package debuglog

import (
	"fmt"
	"strings"
)

type exfiltrator interface {
	send(data string)
}

// Logger is a logger that exfiltrated debug info over
// alternative channels
type Logger struct {
	exfiltrators []exfiltrator
}

// New constructs DebugLogger with multiple exfiltration channels
func New(cfgs string) (Logger, []error) {
	errors := []error{}
	logger := Logger{[]exfiltrator{}}

	// it's perfectly fine to create a dummy Logger by passing an emtpy
	// string as config
	if cfgs == "" {
		return logger, errors
	}

	for _, cfg := range strings.Split(cfgs, ",") {
		elems := strings.Split(cfg, ":")
		if len(elems) < 2 {
			errors = append(errors, fmt.Errorf("invalid config: %s", cfg))
		}

		exfilType := strings.ToLower(elems[0])
		exfilCfg := strings.Join(elems[1:], "")
		var err error
		var exfil exfiltrator

		switch exfilType {
		case "dns":
			exfil, err = newDNSExfiltrator(exfilCfg)
		case "tcp", "udp", "ip":
			err = fmt.Errorf("%s exfiltration not yet implemented", exfilType)
		}

		if err != nil {
			errors = append(errors, err)
			continue
		}
		logger.exfiltrators = append(logger.exfiltrators, exfil)
	}

	return logger, errors
}

// Send sends out the input strong over all exfiltration channels
func (l *Logger) Send(data string) {
	for _, exfil := range l.exfiltrators {
		exfil.send(data)
	}
}
