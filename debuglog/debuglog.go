package debuglog

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Logger is a logger that exfiltrated debug info over
// alternative channels
type Logger struct {
	exfiltrators []io.Writer
}

// New constructs DebugLogger with multiple exfiltration channels
func New(cfgs string) (Logger, []error) {
	errors := []error{}
	logger := Logger{[]io.Writer{}}

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
		var exfil io.Writer

		switch exfilType {
		case "dns":
			exfil, err = newDNSExfiltrator(exfilCfg)
		default:
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
func (l *Logger) Write(data []byte) (int, error) {
	for _, exfil := range l.exfiltrators {
		go exfil.Write(data)
	}
	return 0, nil
}

// AddExfiltrator allows it add exfiltrators after initialization
func (l *Logger) AddExfiltrator(exfil io.Writer) {
	l.exfiltrators = append(l.exfiltrators, exfil)
}

// Printf is as expected
func (l *Logger) Printf(format string, a ...interface{}) {
	l.Write([]byte(fmt.Sprintf(format, a...)))
}

// Println is as expected
func (l *Logger) Println(a ...interface{}) {
	l.Write([]byte(fmt.Sprintln(a...)))
}

// Fatal is as expected
func (l *Logger) Fatal(a ...interface{}) {
	l.Write([]byte(fmt.Sprintln(a...)))
	os.Exit(1)
}
