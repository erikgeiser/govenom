package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type regularLDFlags map[string]string
type externalVarLDFlags map[string]string

func setupLDFlags(regular regularLDFlags, externalVars externalVarLDFlags) string {
	res := ""
	for k, v := range regular {
		res += "-" + k
		if v != "" {
			res += "=" + v
		}
		res += " "
	}
	for k, v := range externalVars {
		res += fmt.Sprintf("-X main.%s=%s ", k, v)
	}
	return res
}

func outputFileName(payload string) string {
	if opts.output != "" {
		return opts.output
	}
	var extension string
	switch opts.os {
	case "windows":
		extension = "exe"
	case "darwin":
		extension = "macho"
	case "linux":
		extension = "elf"
	default:
		extension = "executable"
	}
	return fmt.Sprintf("%s.%s-%s.%s", payload, opts.os, opts.arch, extension)
}

func build(args []string) error {
	a := append([]string{"build"}, args...)

	fmt.Printf("go %s\n\n", strings.Join(a, " "))

	cmd := exec.Command("go", a...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = buildEnv()
	return cmd.Run()
}

func buildEnv() []string {
	envs := []string{}
	for _, env := range os.Environ() {
		if !strings.HasPrefix(env, "GOOS=") && !strings.HasPrefix(env, "GOARCH=") {
			envs = append(envs, env)
		}
	}
	return append(envs, fmt.Sprintf("GOOS=%s", opts.os), fmt.Sprintf("GOARCH=%s", opts.arch))
}
