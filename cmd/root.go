package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/infralight/cli/client"
	"github.com/infralight/cli/config"
	"github.com/infralight/cli/tui"
	"github.com/spf13/cobra"
)

var profile, authHeader, apiURL, accessKey, secretKey string
var c *client.Client
var prettyPrint, failOnError bool

var rootCmd = &cobra.Command{
	Use:   "infralight",
	Short: "Command line interface for the Infralight SaaS",
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		if cmd.Name() == "configure" || cmd.Name() == "version" {
			return nil
		}

		if accessKey == "" || secretKey == "" {
			// keypair not provided via command line flags, try to load a
			// configuration file. If this fails, only exit if user supplied a
			// profile other than default, or if the profile is default and the
			// error was not that the configuration file doesn't exist. This
			// ensures that if no configuration file exists at all, we will
			// prompt the user for a keypair
			conf, err := config.LoadConfig(profile)
			if err != nil {
				if profile == "default" && errors.Is(err, config.ErrConfigNotFound) {
					// no configuration profiles exist yet, so force the user to
					// configure
					profile, err = tui.StartConfigure("No profile exists, please create one")
					if err != nil {
						return err
					}

					// reload this configuration
					conf, err = config.LoadConfig(profile)
				}

				if err != nil {
					return err
				}
			}

			accessKey = conf.AccessKey
			secretKey = conf.SecretKey
			apiURL = conf.URL
			authHeader = conf.AuthorizationHeader
		}

		c = client.New(apiURL, authHeader)

		if accessKey == "" && secretKey == "" {
			return errors.New("access and secret keys must be provided")
		}

		fmt.Fprintf(os.Stderr, "Using profile %q against %q...\n\n", profile, apiURL)

		err := c.Authenticate(accessKey, secretKey)
		if err != nil {
			// in case the new access key / secret are not yet registered in Auth0, we do sleep and retry
			time.Sleep(10 * time.Second)
			return c.Authenticate(accessKey, secretKey)
		}
		return nil
	},
	RunE: func(_ *cobra.Command, _ []string) error {
		return tui.Start(c, accessKey, secretKey)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&profile,
		"profile",
		"p",
		"default",
		"Profile to use",
	)
	rootCmd.PersistentFlags().StringVarP(
		&accessKey,
		"access-key",
		"u",
		"",
		"Access key (will be prompted for if not provided)",
	)
	rootCmd.PersistentFlags().StringVarP(
		&secretKey,
		"secret-key",
		"s",
		"",
		"Secret key (will be prompted for if not provided)",
	)
	rootCmd.PersistentFlags().StringVar(
		&apiURL,
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
	rootCmd.PersistentFlags().BoolVar(
		&failOnError,
		"fail-on-error",
		false,
		"Exit with a non-success code when errors are encountered",
	)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		if failOnError {
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	}
}

func render(v interface{}) error {
	enc := json.NewEncoder(os.Stdout)
	if prettyPrint {
		enc.SetIndent("", "    ")
	}

	return enc.Encode(v)
}
