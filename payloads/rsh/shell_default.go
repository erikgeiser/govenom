// +build !windows

package main

import "syscall"

var suggestedShells = [][]string{
	{"bash", "-i"},
	{"sh"},
}

func sysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{}
}
