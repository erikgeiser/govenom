package exfilwriter

import (
	"os"
)

type fileExfiltrator struct {
	fileName string
}

func newFileExfiltrator(fileName string) (*fileExfiltrator, error) {
	return &fileExfiltrator{fileName}, nil
}

func (ex *fileExfiltrator) Write(data []byte) (int, error) {
	f, err := os.OpenFile(ex.fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return 0, err
	}

	n, err := f.Write(data)
	if err != nil {
		return n, err
	}

	err = f.Close()
	if err != nil {
		return n, err
	}

	return n, nil
}
