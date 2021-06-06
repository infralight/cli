package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var envsCmd = &cobra.Command{
	Use:   "envs [cmd]",
	Short: "Manage Infralight Environments",
	Args:  cobra.MinimumNArgs(1),
}

var envsListCmd = &cobra.Command{
	Use:           "list",
	Short:         "List Infralight Environments",
	Args:          cobra.NoArgs,
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, _ []string) error {
		list, err := c.ListStacks()
		if err != nil {
			return fmt.Errorf("failed listing environments")
		}

		return render(list)
	},
}

func init() {
	envsCmd.AddCommand(envsListCmd)
	rootCmd.AddCommand(envsCmd)
}
