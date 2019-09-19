package exfilwriter

import (
	"fmt"
	"io"
	"strings"
)

// ExfilWriter is a logger that exfiltrated debug info over
// alternative channels
type ExfilWriter struct {
	exfiltrators []io.Writer
}

// New constructs DebugLogger with multiple exfiltration channels
func New(cfgs string) (*ExfilWriter, []error) {
	errors := []error{}
	exfilWriter := ExfilWriter{[]io.Writer{}}

	// it's perfectly fine to create a dummy Logger by passing an emtpy
	// string as config
	if cfgs == "" {
		return nil, errors
	}

	for _, cfg := range strings.Split(cfgs, ",") {
		elems := strings.Split(cfg, ":")
		if len(elems) < 2 {
			errors = append(errors, fmt.Errorf("invalid config: %s", cfg))
		}

		exfilType := strings.ToLower(elems[0])
		exfilCfg := strings.Join(elems[1:], "")
		var err error
		var exfil io.Writer

		switch exfilType {
		case "dns":
			exfil, err = newDNSExfiltrator(exfilCfg)
		case "file":
			exfil, err = newFileExfiltrator(exfilCfg)
		case "dial":
			exfil, err = newDialExfiltrator(exfilCfg)
		default:
			err = fmt.Errorf("%s exfiltration not yet implemented", exfilType)
		}

		if err != nil {
			errors = append(errors, err)
			continue
		}
		exfilWriter.exfiltrators = append(exfilWriter.exfiltrators, exfil)
	}

	return &exfilWriter, errors
}

// Send sends out the input strong over all exfiltration channels
func (l *ExfilWriter) Write(data []byte) (int, error) {
	for _, exfil := range l.exfiltrators {
		go exfil.Write(data)
	}
	return 0, nil
}

// AddExfiltrator allows it add exfiltrators after initialization
func (l *ExfilWriter) AddExfiltrator(exfil io.Writer) {
	l.exfiltrators = append(l.exfiltrators, exfil)
}
