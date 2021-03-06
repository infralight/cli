package cmd

import (
	"fmt"
	"strings"

	"github.com/infralight/cli/version"
	"github.com/spf13/cobra"
)

var (
	environmentCreateType   string
	environmentCreateName   string
	environmentCreateOwner  string
	environmentCreateLabels []string
	environmentDeleteId     string
)

var envsCmd = &cobra.Command{
	Use:   "envs [cmd]",
	Short: fmt.Sprintf("Manage %s Environments", version.Product),
	Args:  cobra.MinimumNArgs(1),
}

var envsListCmd = &cobra.Command{
	Use:           "list",
	Short:         fmt.Sprintf("List %s Environments", version.Product),
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
	Use:           "delete --env-id ENVIRONMENT_ID",
	Short:         fmt.Sprintf("Delete %s Environment", version.Product),
	Args:          cobra.NoArgs,
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, args []string) error {
		env, err := c.DeleteEnvironment(environmentDeleteId)
		if err != nil {
			return fmt.Errorf("failed deleting environment")
		}

		return render(env)
	},
}

// nolint:lll
var envCreateCmd = &cobra.Command{
	Use:           "create --name ENVIRONMENT_NAME --type ENVIRONMENT_TYPE --owner ENVIRONMENT_OWNER --labels=label1,label2",
	Short:         fmt.Sprintf("Create %s Environment", version.Product),
	Args:          cobra.NoArgs,
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, args []string) error {
		env, err := c.CreateEnvironment(environmentCreateName, environmentCreateType, environmentCreateOwner, environmentCreateLabels)
		if err != nil {
			return fmt.Errorf("failed create environment")
		}

		return render(env)
	},
}

//nolint
func init() {
	envCreateCmd.PersistentFlags().StringVar(&environmentCreateName, "name", "", "Environment Name")
	envCreateCmd.MarkPersistentFlagRequired("name")
	envCreateCmd.PersistentFlags().StringVar(&environmentCreateType, "type", "iacStack", "Environment Type")
	envCreateCmd.PersistentFlags().StringVar(&environmentCreateOwner, "owner", fmt.Sprintf("%sCli", version.Product), fmt.Sprintf("Environment owner. Use `%s users list` to see available options", strings.ToLower(version.Product)))
	envCreateCmd.MarkPersistentFlagRequired("owner")
	envCreateCmd.PersistentFlags().StringSliceVar(&environmentCreateLabels, "labels", nil, "Environment labels.")
	envDeleteCmd.PersistentFlags().StringVar(&environmentDeleteId, "env-id", "", "The environment id to list the stacks for")
	envDeleteCmd.MarkPersistentFlagRequired("env-id")
	envsCmd.AddCommand(envCreateCmd)
	envsCmd.AddCommand(envsListCmd)
	envsCmd.AddCommand(envDeleteCmd)
	rootCmd.AddCommand(envsCmd)
}
