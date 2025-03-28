package cmd

import (
	"fmt"

	"github.com/metal-stack/v"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show currently used version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s\n", v.V.String())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
