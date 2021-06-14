package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var statePoliciesCmd = &cobra.Command{
	Use:   "state-policies [cmd]",
	Short: "Manage State Policies",
	Args:  cobra.MinimumNArgs(1),
}

var statePoliciesGetCmd = &cobra.Command{
	Use:           "get <stack>",
	Short:         "Get State Policy In a Stack",
	Args:          cobra.ExactArgs(1),
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, args []string) error {
		policy, err := c.GetStatePolicy(args[0])
		if err != nil {
			return fmt.Errorf("failed getting state policy: %w", err)
		}

		return render(policy)
	},
}

var statePoliciesUpdateCmd = &cobra.Command{
	Use:           "update <stack> <selected> <policy.json>",
	Short:         "Update a Stack's State Policy",
	Args:          cobra.ExactArgs(3),
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, args []string) error {
		policy, err := os.ReadFile(args[2])
		if err != nil {
			return fmt.Errorf("failed reading %s: %w", args[2], err)
		}

		err = c.UpdateStatePolicy(args[0], args[1], policy)
		if err != nil {
			return fmt.Errorf("failed updating state policy: %w", err)
		}

		print("Success\n")
		return nil
	},
}

func init() {
	statePoliciesCmd.AddCommand(statePoliciesGetCmd, statePoliciesUpdateCmd)
	rootCmd.AddCommand(statePoliciesCmd)
}
