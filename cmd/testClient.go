package cmd

import (
	"os"

	"github.com/insomniacslk/dhcp/dhcpv4/client4"
	"github.com/spf13/cobra"
)

var testClientCmd = &cobra.Command{
	Use:   "test-client",
	Short: "start a fake client to test relay functionality",
	Run: func(cmd *cobra.Command, args []string) {
		intf, _ := cmd.Flags().GetString("interface")

		client := client4.NewClient()
		conversation, err := client.Exchange(intf)

		// print the packets before handling error because Exchange can return a non-empty packet list even if there was an error
		for _, packet := range conversation {
			log.Debug("packet received", "packet", packet.Summary())
		}

		if err != nil {
			log.Error("error occurred during packet exchange", "error", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(testClientCmd)
	testClientCmd.Flags().StringP("interface", "i", "eth0", "interface to configure via DHCP")
}
