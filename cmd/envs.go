package cmd

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/infralight/cli/client"
	"github.com/spf13/cobra"
)

var envsCmd = &cobra.Command{
	Use:   "envs [cmd]",
	Short: "Manage Infralight Environments",
	Args:  cobra.MinimumNArgs(1),
	PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
		if accessKey == "" || secretKey == "" {
			return errors.New("access and secret keys must be provided")
		}

		c = client.New(url, authHeader)
		return c.Authenticate(accessKey, secretKey)
	},
}

var envsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Infralight Environments",
	Args:  cobra.NoArgs,
	RunE: func(_ *cobra.Command, _ []string) error {
		list, err := c.ListStacks()
		if err != nil {
			return err
		}

		enc := json.NewEncoder(os.Stdout)
		if prettyPrint {
			enc.SetIndent("", "    ")
		}

		return enc.Encode(map[string]interface{}{
			"environments": list,
		})
	},
}

func init() {
	envsCmd.AddCommand(envsListCmd)
	rootCmd.AddCommand(envsCmd)
}
