package cmd

import (
	"govenom/tools/dnslogger"
	"govenom/tools/gateway"
	"govenom/tools/pusher"

	"github.com/spf13/cobra"
)

var toolCmd = &cobra.Command{
	Use:           "tool",
	Short:         "run a tool",
	SilenceErrors: true,
}

func init() {
	toolCmd.AddCommand(dnslogger.DNSLoggerCmd)
	toolCmd.AddCommand(pusher.PusherCmd)
	toolCmd.AddCommand(gateway.GatewayCmd)
}
