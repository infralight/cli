package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	environmentCreateType string
	environmentCreateName string
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

var envDeleteCmd = &cobra.Command{
	Use:           "delete",
	Short:         "Delete Infralight Environment",
	Args:          cobra.ExactArgs(1),
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, args []string) error {
		env, err := c.DeleteEnvironment(args[0])
		if err != nil {
			return fmt.Errorf("failed deleting environment")
		}

		return render(env)
	},
}

var envCreateCmd = &cobra.Command{
	Use:           "create ENVIRONMENT_NAME ENVIRONMENT_TYPE",
	Short:         "Create Infralight Environment",
	Args:          cobra.NoArgs,
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, args []string) error {
		env, err := c.CreateEnvironment(environmentCreateName, environmentCreateType)
		if err != nil {
			return fmt.Errorf("failed create environment")
		}

		return render(env)
	},
}

func init() {
	envCreateCmd.PersistentFlags().StringVar(&environmentCreateName, "name", "", "Environment Name")
	envCreateCmd.PersistentFlags().StringVar(&environmentCreateType, "type", "iacStack", "Environment Type")
	envsCmd.AddCommand(envCreateCmd)
	envsCmd.AddCommand(envsListCmd)
	envsCmd.AddCommand(envDeleteCmd)
	rootCmd.AddCommand(envsCmd)
}
