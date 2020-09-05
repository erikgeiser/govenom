package dnslogger

import (
	"fmt"
	"log"

	"github.com/miekg/dns"
	"github.com/spf13/cobra"
)

var opts struct {
	net     string
	address string
	verbose bool
}

// DNSLoggerCmd contains the CLI interface for the dnslogger command.
var DNSLoggerCmd = &cobra.Command{
	Use:   "dnslogger",
	Short: "dns logger is a receiver for DNS-exfiltrated logs of govenom payloads",
	Run: func(cmd *cobra.Command, args []string) {
		handler := newDNSHandler(newSimpleLogHandler())
		srv := &dns.Server{Addr: opts.address, Net: opts.net, Handler: handler}

		fmt.Printf("Setting up listener on %s/%s\n", srv.Addr, srv.Net)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Failed to set udp listener %s\n", err.Error())
		}
	},
}

func init() {
	DNSLoggerCmd.PersistentFlags().StringVarP(&opts.net, "net", "n", "udp", "network protocol")
	DNSLoggerCmd.PersistentFlags().StringVarP(&opts.address, "address", "a", ":53", "listen address ([ip]:port)")
	DNSLoggerCmd.PersistentFlags().BoolVarP(&opts.verbose, "verbose", "v", false, "verbose error logging")
}

func logError(err error) {
	if !opts.verbose {
		return
	}

	fmt.Printf("%s\n", err)
}
