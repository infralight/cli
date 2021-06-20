package cmd

import (
	"fmt"
	"os"

	"github.com/infralight/cli/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the CLI version and exist",
	Args:  cobra.NoArgs,
	RunE: func(_ *cobra.Command, _ []string) error {
		_, err := fmt.Fprintln(os.Stdout, version.Version)
		return err
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
