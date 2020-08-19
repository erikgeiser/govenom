// +build !windows

package main

import (
	"fmt"
	"path"
	"syscall"
)

func buildShellCommandList(prioritizedChoices ...[]string) ([][]string, error) {
	shells := prioritizedChoices
	shells = append(shells, [][]string{
		{"bash", "-i"},
		{"sh"},
		{"zsh"},
		{"csh"},
		{"dash"},
		{"ash"},
	}...)

	prefixes := []string{
		"",
		"/bin",
		"/sbin",
		"/usr/bin",
		"/usr/sbin",
		"/usr/local/bin",
		"/usr/local/sbin",
	}

	cmds := make([][]string, 0, len(shells)*len(prefixes))

	for _, shell := range shells {
		for _, prefix := range prefixes {
			if len(shell) == 0 {
				continue
			}

			if shell[0] == "" {
				continue
			}

			cmd := []string{path.Join(prefix, shell[0])}

			cmds = append(cmds, append(cmd, shell[1:]...))
		}
	}

	if len(cmds) == 0 {
		return nil, fmt.Errorf("no suggested shells")
	}

	return cmds, nil
}

func sysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{}
}
