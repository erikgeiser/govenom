package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/miekg/dns"
	"github.com/spf13/cobra"
)

var opts struct {
	multiplexingEnabled  bool
	dnsMessageIDLength   int
	interFragmentTimeout time.Duration
	net                  string
	address              string
	verbose              bool
}

var rootCmd = &cobra.Command{
	Use:   "dns_logger",
	Short: "dns logger is a receiver for DNS-exfiltrated logs of govenom payloads",
	Run: func(cmd *cobra.Command, args []string) {
		var logger logHandler
		if opts.multiplexingEnabled {
			logger = newMultiplexingLogHandler(opts.dnsMessageIDLength, opts.interFragmentTimeout)
		} else {
			logger = newSimpleLogHandler()
		}

		handler := newDNSHandler(logger)
		srv := &dns.Server{Addr: opts.address, Net: opts.net, Handler: handler}

		fmt.Printf("Setting up listener on %s/%s\n", srv.Addr, srv.Net)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Failed to set udp listener %s\n", err.Error())
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&opts.net, "net", "n", "udp", "network protocol")
	rootCmd.PersistentFlags().StringVarP(&opts.address, "address", "a", ":53", "listen address ([ip]:port)")
	rootCmd.PersistentFlags().BoolVarP(&opts.verbose, "verbose", "v", false, "verbose error logging")
	rootCmd.PersistentFlags().BoolVarP(&opts.multiplexingEnabled, "multiplexing", "m", false, "enable message multiplexing")
	rootCmd.PersistentFlags().IntVarP(&opts.dnsMessageIDLength, "id-length", "l", 6, "multiplexing message ID length")
	rootCmd.PersistentFlags().DurationVarP(&opts.interFragmentTimeout, "timeout", "t", 500*time.Microsecond, "multiplexing inter-fragment-timeout")
}

func logError(err error) {
	if !opts.verbose {
		return
	}
	fmt.Printf("%s\n", err)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
