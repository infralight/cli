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
		list, err := c.ListEnvironments()
		if err != nil {
			return fmt.Errorf("failed listing environments")
		}

		return render(list)
	},
}

var envCreateCmd = &cobra.Command{
	Use:           "create ENVIRONMENT_NAME ENVIRONMENT_TYPE",
	Short:         "Create Infralight Environment",
	Args:          cobra.ExactArgs(2),
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, args []string) error {
		env, err := c.CreateEnvironment(args[0], args[1])
		if err != nil {
			return fmt.Errorf("failed create environment")
		}

		return render(env)
	},
}

func init() {
	envsCmd.AddCommand(envsListCmd)
	rootCmd.AddCommand(envCreateCmd)
	rootCmd.AddCommand(envsCmd)
}
