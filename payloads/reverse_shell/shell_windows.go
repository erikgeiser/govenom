// +build windows

package main

import (
	"fmt"
	"os"
	"syscall"
)

func buildShellCommandList(prioritizedChoices ...[]string) ([][]string, error) {
	windir := os.Getenv("windir")
	if windir == "" {
		windir = "C:\\Windows"
	}

	cmds := prioritizedChoices
	cmds = append(cmds, [][]string{
		{"powershell"},
		{fmt.Sprintf("%s\\system32\\WindowsPowerShell\\v1.0\\powershell.exe", windir)}, // PS x64
		{fmt.Sprintf("%s\\syswow64\\WindowsPowerShell\\v1.0\\powershell.exe", windir)}, // PS x32
		{"cmd"},
		{fmt.Sprintf("%s\\system32\\cmd.exe", windir)}, // cmd

	}...)

	return cmds, nil
}

func getSysProcAttr() *syscall.SysProcAttr {
	if noWindowsGui == "true" {
		return &syscall.SysProcAttr{}
	}

	return &syscall.SysProcAttr{HideWindow: true}
}
