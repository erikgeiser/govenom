// +build windows

package main

import (
	"fmt"
	"os"
	"syscall"
)

func getShellBinaries() []string {
	windir := os.Getenv("windir")
	if windir == "" {
		windir = "C:\\Windows"
	}

	return []string{
		"powershell",
		"bash",
		fmt.Sprintf("%s\\system32\\WindowsPowerShell\\v1.0\\powershell.exe", windir), // PS x64
		fmt.Sprintf("%s\\syswow64\\WindowsPowerShell\v1.0\\powershell.exe", windir),  // PS x32
		"cmd",
		fmt.Sprintf("%s\\system32\\cmd.exe", windir), // cmd
	}
}

func getSysProcAttr() *syscall.SysProcAttr {
	if noWindowsGui == "true" {
		return &syscall.SysProcAttr{}
	}
	return &syscall.SysProcAttr{HideWindow: true}
}
