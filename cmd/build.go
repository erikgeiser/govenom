package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/markbates/pkger"
)

type injectedVariables map[string]string

func setupLDFlags(vars injectedVariables, opts BuildOpts) string {
	ldFlags := "-w -s "

	if opts.OS == "windows" && !opts.NoWindowsGui {
		ldFlags += "-H=windowsgui"
	}

	for k, v := range vars {
		ldFlags += fmt.Sprintf(`-X "main.%s=%s" `, k, v)
	}

	return ldFlags
}

func outputFileName(payload string, opts BuildOpts) (string, error) {
	if opts.Output != "" {
		return filepath.Abs(opts.Output)
	}

	var extension string

	switch opts.OS {
	case "windows":
		extension = "exe"
	case "darwin":
		extension = "macho"
	case "linux":
		extension = "elf"
	default:
		extension = "executable"
	}

	fileName := fmt.Sprintf("%s.%s-%s.%s", payload, opts.OS, opts.Arch, extension)

	return filepath.Abs(fileName)
}

func build(payload string, vars injectedVariables, opts BuildOpts) error {
	_, err := exec.LookPath(opts.GoBin)
	if err != nil {
		return fmt.Errorf("cannot find Go binary (install Go or set `--go /path/to/go`")
	}

	buildDir, err := ioutil.TempDir("", "govenom_build_")
	if err != nil {
		return fmt.Errorf("creating temporary build directory: %v", err)
	}

	defer func() {
		err := os.RemoveAll(buildDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: removing build dir %q: %v\n", buildDir, err)
		}
	}()

	err = pkger.Walk("/payloads", pkgerCopyWalker(path.Join(buildDir, "payloads")))
	if err != nil {
		return fmt.Errorf("extracting payload: %w", err)
	}

	err = ioutil.WriteFile(path.Join(buildDir, "go.mod"),
		[]byte("module govenom\n\n\ngo 1.11"), 0600)
	if err != nil {
		return fmt.Errorf("extracting go.mod")
	}

	outFileName, err := outputFileName(payload, opts)
	if err != nil {
		return fmt.Errorf("determine absolute output file name: %w", err)
	}

	args := []string{
		"go", "build",
		"-trimpath",
		"-ldflags", setupLDFlags(vars, opts),
		"-o", outFileName,
		"./payloads/" + payload,
	}

	fmt.Printf("Compiling: [\"%s\"]\n", strings.Join(args, "\", \""))

	cmd := exec.Command(args[0], args[1:]...) // nolint:gosec
	cmd.Dir = buildDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = buildEnv(opts)

	return cmd.Run()
}

func pkgerCopyWalker(dst string) filepath.WalkFunc {
	return func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		srcFilePath := filePath
		dstFilePath := path.Join(dst,
			strings.TrimPrefix(srcFilePath, "govenom:/payloads/"))

		if info.IsDir() {
			return os.MkdirAll(dstFilePath, os.ModePerm)
		}

		return pkgerCopyFile(dstFilePath, srcFilePath)
	}
}

func pkgerCopyFile(dst string, pkgerFilePath string) error {
	srcFile, err := pkger.Open(pkgerFilePath)
	if err != nil {
		return fmt.Errorf("open embedded file: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("create dst file: %w", err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("copying embedded file: %w", err)
	}

	err = srcFile.Close()
	if err != nil {
		return fmt.Errorf("closing embedded file: %w", err)
	}

	err = dstFile.Close()
	if err != nil {
		return fmt.Errorf("closing dst file: %w", err)
	}

	return nil
}

func buildEnv(opts BuildOpts) []string {
	envs := []string{}

	for _, env := range os.Environ() {
		if !strings.HasPrefix(env, "GOOS=") && !strings.HasPrefix(env, "GOARCH=") {
			envs = append(envs, env)
		}
	}

	return append(envs, fmt.Sprintf("GOOS=%s", opts.OS), fmt.Sprintf("GOARCH=%s", opts.Arch))
}
