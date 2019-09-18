// +build !windows

package main

import (
	"path"
	"syscall"
)

func getShellBinaries() []string {
	shells := []string{
		"bash",
		"sh",
		"zsh",
		"csh",
		"dash",
		"ash",
	}
	prefixes := []string{
		"",
		"/bin",
		"/sbin",
		"/usr/bin",
		"/usr/sbin",
		"/usr/local/bin",
		"/usr/local/sbin",
	}
	binaries := make([]string, 0, len(shells)*len(prefixes))
	for _, prefix := range prefixes {
		for _, shell := range shells {
			binaries = append(binaries, path.Join(prefix, shell))
		}
	}
	return binaries
}

func getSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{}
}
