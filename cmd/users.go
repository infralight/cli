package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var usersCmd = &cobra.Command{
	Use:   "users [cmd]",
	Short: "Manage Users",
	Args:  cobra.MinimumNArgs(1),
}

var usersListCmd = &cobra.Command{
	Use:           "list",
	Short:         "List users",
	Args:          cobra.NoArgs,
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, args []string) error {
		list, err := c.ListUsers()
		if err != nil {
			return fmt.Errorf("failed listing states: %w", err)
		}

		return render(list)
	},
}

func init() {
	usersCmd.AddCommand(usersListCmd)
	rootCmd.AddCommand(usersCmd)
}
