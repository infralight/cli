package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var stacksCmd = &cobra.Command{
	Use:   "stacks [cmd]",
	Short: "Manage Stacks",
	Args:  cobra.MinimumNArgs(1),
}

var stacksListCmd = &cobra.Command{
	Use:           "list <env_id>",
	Short:         "List States in an Environment",
	Args:          cobra.ExactArgs(1),
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, args []string) error {
		list, err := c.ListStacks(args[0])
		if err != nil {
			return fmt.Errorf("failed listing stacks: %w", err)
		}

		return render(list)
	},
}

var stacksGetCmd = &cobra.Command{
	Use:           "get <environment_id> <stack_id>",
	Short:         "Get a Specific Stack",
	Args:          cobra.ExactArgs(2),
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, args []string) error {
		stack, err := c.GetStack(args[0], args[1])
		if err != nil {
			return fmt.Errorf("failed getting stack %s:%s: %w", args[0], args[1], err)
		}

		return render(stack)
	},
}

var stacksDeleteCmd = &cobra.Command{
	Use:           "delete <environment_id> <stack_id>",
	Short:         "Delete a Specific Stack",
	Args:          cobra.ExactArgs(2),
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, args []string) error {
		stack, err := c.DeleteStack(args[0], args[1])
		if err != nil {
			return fmt.Errorf("failed deleting stack %s:%s: %w", args[0], args[1], err)
		}

		return render(stack)
	},
}

func init() {
	stacksCmd.AddCommand(stacksListCmd, stacksGetCmd, stacksDeleteCmd)
	rootCmd.AddCommand(stacksCmd)
}
