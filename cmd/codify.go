package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var codifyCmd = &cobra.Command{
	Use:   "codify <provider> <provider_id> <asset_type> <asset_id>",
	Short: "Codify an unmanaged asset",
	Args:  cobra.ExactArgs(4),
	RunE: func(_ *cobra.Command, args []string) error {
		output, err := c.Codify(args[0], args[1], args[2], args[3])
		if err != nil {
			return err
		}

		_, err = fmt.Fprintln(os.Stdout, output)
		return err
	},
}

func init() {
	rootCmd.AddCommand(codifyCmd)
}
