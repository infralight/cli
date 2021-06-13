package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var statesCmd = &cobra.Command{
	Use:   "states [cmd]",
	Short: "View States",
	Args:  cobra.MinimumNArgs(1),
}

var statesListCmd = &cobra.Command{
	Use:           "list <stack>",
	Short:         "List States In a Stack",
	Args:          cobra.ExactArgs(1),
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, args []string) error {
		list, err := c.ListStates(args[0])
		if err != nil {
			return fmt.Errorf("failed listing states: %w", err)
		}

		return render(list)
	},
}

var statesLatestCmd = &cobra.Command{
	Use:           "latest <stack>",
	Short:         "Get Latest State In a Stack",
	Args:          cobra.ExactArgs(1),
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, args []string) error {
		list, err := c.GetLatestState(args[0])
		if err != nil {
			return fmt.Errorf("failed getting latest state: %w", err)
		}

		return render(list)
	},
}

var statesUploadCmd = &cobra.Command{
	Use:           "upload <stack> <tf_state.json> <policy.json>",
	Short:         "Upload a Policy for a State",
	Args:          cobra.ExactArgs(3),
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, args []string) error {
		tfState, err := os.ReadFile(args[1])
		if err != nil {
			return fmt.Errorf("failed reading %s: %w", args[1], err)
		}

		policy, err := os.ReadFile(args[2])
		if err != nil {
			return fmt.Errorf("failed reading %s: %w", args[2], err)
		}

		err = c.UploadStatePolicy(args[0], tfState, policy)
		if err != nil {
			return fmt.Errorf("failed uploading policy: %w", err)
		}

		print("Success\n")
		return nil
	},
}

func init() {
	statesCmd.AddCommand(statesListCmd, statesLatestCmd, statesUploadCmd)
	rootCmd.AddCommand(statesCmd)
}
