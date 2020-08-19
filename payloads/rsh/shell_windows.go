// +build windows

package main

import "syscall"

var suggestedShells = [][]string{
	{"powershell"},
	{"cmd"},
}

func sysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{HideWindow: !noWindowsGuiValue}
}
