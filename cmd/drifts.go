package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var driftsCmd = &cobra.Command{
	Use:   "drifts [cmd]",
	Short: "View Calculated Drifts",
	Args:  cobra.MinimumNArgs(1),
}

var (
	driftsListOnlyDelta bool
	driftsListLimit     uint64
)

var driftsListCmd = &cobra.Command{
	Use:           "list",
	Short:         "List Drifts",
	Args:          cobra.NoArgs,
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, _ []string) error {
		list, err := c.ListDrifts(driftsListOnlyDelta, driftsListLimit)
		if err != nil {
			return fmt.Errorf("failed listing drifts: %w", err)
		}

		return render(list)
	},
}

var driftsShowCmd = &cobra.Command{
	Use:           "show <drift_id>",
	Short:         "Show a Specific Drift",
	Args:          cobra.ExactArgs(1),
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, args []string) error {
		driftID := strings.TrimPrefix(args[0], "Drifts/")

		drift, err := c.ShowDrift(driftID)
		if err != nil {
			return fmt.Errorf("failed showing drift %s: %w", driftID, err)
		}

		return render(drift)
	},
}

var driftsAssetCmd = &cobra.Command{
	Use:           "asset <asset_id>",
	Short:         "Show a Specific Asset",
	Args:          cobra.ExactArgs(1),
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, args []string) error {
		asset, err := c.ShowAsset(args[0])
		if err != nil {
			return fmt.Errorf("failed showing asset %s: %w", args[0], err)
		}

		return render(asset)
	},
}

func init() {
	driftsListCmd.Flags().BoolVar(&driftsListOnlyDelta, "only-delta", false, "Only show assets with a detected drift")
	driftsListCmd.Flags().Uint64Var(&driftsListLimit, "limit", 0, "Only show assets with a detected drift")
	driftsCmd.AddCommand(driftsListCmd, driftsShowCmd, driftsAssetCmd)
	rootCmd.AddCommand(driftsCmd)
}
