package cmd

import (
	"fmt"
	"os"

    "github.com/spf13/cobra"
	
	"github.com/infralight/cli/tui"
	"github.com/infralight/cli/version"
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: fmt.Sprintf("Configure %s authentication", version.Product),
	Args:  cobra.NoArgs,
	RunE: func(_ *cobra.Command, _ []string) error {
		profile, err := tui.StartConfigure("")
		if err != nil {
			return err
		}

		fmt.Fprintf(os.Stdout, "Successfully created profile %q.\n", profile)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)
}
