package exfilwriter

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

// ExfilWriter is a logger that exfiltrated debug info over
// alternative channels
type ExfilWriter struct {
	exfiltrators []io.Writer
	writeTimeout time.Duration
}

// New constructs DebugLogger with multiple exfiltration channels
func New(cfgs string, writeTimeout time.Duration) (*ExfilWriter, []error) {
	errors := []error{}
	exfilWriter := ExfilWriter{[]io.Writer{}, writeTimeout}

	// it's perfectly fine to create a dummy Logger
	// by passing an emtpy string as config
	if cfgs == "" {
		return &exfilWriter, errors
	}

	for _, cfg := range strings.Split(cfgs, ",") {
		elems := strings.Split(cfg, ":")
		exfilType := strings.ToLower(elems[0])
		exfilCfg := strings.Join(elems[1:], "")

		var err error
		var exfil io.Writer

		switch exfilType {
		case "stdout":
			exfil, err = newWriterExfiltrator(os.Stdout)
		case "stderr":
			exfil, err = newWriterExfiltrator(os.Stderr)
		case "dns":
			exfil, err = newDNSExfiltrator(exfilCfg)
		case "file":
			exfil, err = newFileExfiltrator(exfilCfg)
		case "dial":
			exfil, err = newDialExfiltrator(exfilCfg)
		default:
			err = fmt.Errorf("%s exfiltration not implemented", exfilType)
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
	var wg sync.WaitGroup

	for _, exfil := range l.exfiltrators {
		wg.Add(1)

		go func(ew io.Writer) {
			defer wg.Done()
			_, _ = ew.Write(data)
		}(exfil)
	}

	waitTimeout(&wg, l.writeTimeout)

	return len(data), nil
}

// AddExfiltrator allows it add exfiltrators after initialization
func (l *ExfilWriter) AddExfiltrator(exfil io.Writer) {
	l.exfiltrators = append(l.exfiltrators, exfil)
}

func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) {
	c := make(chan struct{})

	go func() {
		defer close(c)
		wg.Wait()
	}()

	select {
	case <-c:
	case <-time.After(timeout):
	}
}
