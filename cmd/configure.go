package cmd

import (
	"github.com/infralight/cli/tui"
	"github.com/spf13/cobra"
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure Infralight authentication",
	Args:  cobra.NoArgs,
	RunE: func(_ *cobra.Command, _ []string) error {
		return tui.StartConfigure()
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)
}
