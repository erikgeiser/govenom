package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "govenom",
	Short:         "govenom is a cross-platform payload generator",
	SilenceErrors: true,
}

func init() {
	rootCmd.AddCommand(payloadCmd)
	rootCmd.AddCommand(toolCmd)
}

// Execute starts govemon
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
