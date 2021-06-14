package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var statesCmd = &cobra.Command{
	Use:   "states [cmd]",
	Short: "Manage States",
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
	Use:           "upload <stack> <tf_state.json> [policy.json]",
	Short:         "Upload a State File",
	Args:          cobra.MinimumNArgs(2),
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, args []string) error {
		tfState, err := os.ReadFile(args[1])
		if err != nil {
			return fmt.Errorf("failed reading %s: %w", args[1], err)
		}

		var policy json.RawMessage
		if len(args) > 2 {
			policy, err = os.ReadFile(args[2])
			if err != nil {
				return fmt.Errorf("failed reading %s: %w", args[2], err)
			}
		}

		err = c.UploadState(args[0], tfState, policy)
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
