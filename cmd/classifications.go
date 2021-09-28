package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/infralight/cli/client"
	"github.com/spf13/cobra"
)

var classCmd = &cobra.Command{
	Use:     "classifications [cmd]",
	Short:   "Manage Asset Classifications",
	Args:    cobra.MinimumNArgs(1),
	Aliases: []string{"class"},
}

var classListCmd = &cobra.Command{
	Use:           "list",
	Short:         "List Classifications",
	Args:          cobra.NoArgs,
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, _ []string) error {
		list, err := c.Classifications.List()
		if err != nil {
			return fmt.Errorf("failed listing classifications: %w", err)
		}

		return render(list)
	},
}

var (
	classCreateInput client.CreateClassificationInput
	classCreateCmd   = &cobra.Command{
		Use:           "create",
		Short:         "Create Classification",
		Args:          cobra.NoArgs,
		SilenceErrors: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			if classCreateInput.Rego == "-" {
				// read Rego code from standard input
				b, err := io.ReadAll(os.Stdin)
				if err != nil {
					return fmt.Errorf("failed reading policy from stdin: %w", err)
				}

				classCreateInput.Rego = string(b)
			}

			class, err := c.Classifications.Create(classCreateInput)
			if err != nil {
				return fmt.Errorf("failed creating classification: %w", err)
			}

			return render(class)
		},
	}
)

var (
	classUpdateInput client.UpdateClassificationInput
	classUpdateCmd   = &cobra.Command{
		Use:           "update CLASSIFICATION_ID",
		Short:         "Update Existing Classification",
		Args:          cobra.ExactArgs(1),
		SilenceErrors: true,
		RunE: func(_ *cobra.Command, args []string) error {
			if classUpdateInput.Rego == "-" {
				// read Rego code from standard input
				b, err := io.ReadAll(os.Stdin)
				if err != nil {
					return fmt.Errorf("failed reading policy from stdin: %w", err)
				}

				classUpdateInput.Rego = string(b)
			}

			class, err := c.Classifications.Update(args[0], classUpdateInput)
			if err != nil {
				return fmt.Errorf(
					"failed updating classification %q: %w",
					args[0],
					err,
				)
			}

			return render(class)
		},
	}
)

var classDeleteCmd = &cobra.Command{
	Use:           "delete CLASSIFICATION_ID",
	Short:         "Delete Existing Classification",
	Args:          cobra.ExactArgs(1),
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, args []string) error {
		err := c.Classifications.Delete(args[0])
		if err != nil {
			return fmt.Errorf(
				"failed deleting classification %q: %w",
				args[0],
				err,
			)
		}

		return nil
	},
}

func init() {
	// Classification creation flags
	classCreateCmd.PersistentFlags().StringVar(
		&classCreateInput.Name, "name", "", "Classification name (required)",
	)
	classCreateCmd.MarkPersistentFlagRequired("name") // nolint: errcheck
	classCreateCmd.PersistentFlags().StringVar(
		&classCreateInput.Description,
		"description",
		"",
		"Classification description",
	)
	classCreateCmd.PersistentFlags().StringVar(
		&classCreateInput.Type,
		"type",
		"",
		"Classification type (required)",
	)
	classCreateCmd.MarkPersistentFlagRequired("type") // nolint: errcheck
	classCreateCmd.PersistentFlags().StringSliceVar(
		&classCreateInput.Labels, "labels", nil, "Classification labels",
	)
	classCreateCmd.PersistentFlags().StringVar(
		&classCreateInput.Rego,
		"rego",
		"",
		"Rego policy encoded in Base64 format (required). Use a dash ('-') to read code from standard input",
	)
	classCreateCmd.MarkPersistentFlagRequired("rego") // nolint: errcheck

	// Classification modification flags
	classUpdateCmd.PersistentFlags().StringVar(
		&classUpdateInput.Name, "name", "", "Classification name",
	)
	classUpdateCmd.PersistentFlags().StringVar(
		&classUpdateInput.Description,
		"description",
		"",
		"Classification description",
	)
	classUpdateCmd.PersistentFlags().StringVar(
		&classUpdateInput.Type,
		"type",
		"",
		"Classification type",
	)
	classUpdateCmd.PersistentFlags().StringSliceVar(
		&classUpdateInput.Labels, "labels", nil, "Classification labels",
	)
	classUpdateCmd.PersistentFlags().StringVar(
		&classUpdateInput.Rego,
		"rego",
		"",
		"Rego policy encoded in Base64 format. Use a dash ('-') to read code from standard input",
	)

	classCmd.AddCommand(
		classListCmd,
		classCreateCmd,
		classUpdateCmd,
		classDeleteCmd,
	)
	rootCmd.AddCommand(classCmd)
}
