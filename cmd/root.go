package cmd

import (
	"fmt"
	"os"

	"github.com/infralight/cli/client"
	"github.com/infralight/cli/tui"
	"github.com/spf13/cobra"
)

var authHeader, url, accessKey, secretKey string
var c *client.Client
var prettyPrint bool

var rootCmd = &cobra.Command{
	Use:   "infralight",
	Short: "Command line interface for the Infralight SaaS",
	RunE: func(_ *cobra.Command, _ []string) error {
		return tui.Start(client.New(url, authHeader), accessKey, secretKey)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&accessKey,
		"access-key",
		"u",
		"",
		"access key (will be prompted for if not provided)",
	)
	rootCmd.PersistentFlags().StringVarP(
		&secretKey,
		"secret-key",
		"p",
		"",
		"secret key (will be prompted for if not provided)",
	)
	rootCmd.PersistentFlags().StringVar(
		&url,
		"url",
		client.DefaultInfralightURL,
		"Infralight API URL",
	)
	rootCmd.PersistentFlags().StringVar(
		&authHeader,
		"auth-header",
		client.DefaultAuthHeader,
		"Authorization header",
	)
	rootCmd.PersistentFlags().BoolVar(
		&prettyPrint,
		"pretty",
		false,
		"Pretty-print JSON output",
	)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
