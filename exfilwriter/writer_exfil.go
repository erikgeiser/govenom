package exfilwriter

import "io"

type writerExfiltrator struct {
	io.Writer
}

func newWriterExfiltrator(w io.Writer) (*writerExfiltrator, error) {
	return &writerExfiltrator{w}, nil
}
