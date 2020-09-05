package gateway

import (
	"log"

	"github.com/spf13/cobra"
)

// Opts hold the options for the gateway command
type Opts struct {
	connectBackAddress string
	gatewayAddress     string
}

var gatewayOpts Opts

// GatewayCmd contains the CLI interface for the gateway command.
var GatewayCmd = &cobra.Command{
	Use:   "gateway",
	Short: "A gateway for the reverse socks5 payload",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runGateway(gatewayOpts)
	},
}

func init() {
	flags := GatewayCmd.PersistentFlags()
	flags.StringVarP(&gatewayOpts.connectBackAddress, "connectback", "c", "",
		"Connect back `listen address` to which the payload will connect")
	flags.StringVarP(&gatewayOpts.gatewayAddress, "gateway", "g", "",
		"Gateway `listen address` which exposes the payloads socks5 proxy")

	_ = GatewayCmd.MarkPersistentFlagRequired("connectback")
	_ = GatewayCmd.MarkPersistentFlagRequired("gateway")
}

func runGateway(opts Opts) error {
	for {
		err := runGatwayForPayloadConnection(opts)
		if err != nil {
			log.Printf("gateway error: %v", err)
		}
	}
}
