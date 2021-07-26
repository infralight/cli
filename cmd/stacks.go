package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	stacksCreateLocation         string
	stacksCreateStatePath        string
	stacksCreateS3IntegrationID  string
	stacksCreateS3Bucket         string
	stacksCreateS3Key            string
	stacksCreateTfCAPIToken      string
	stacksCreateTfCWorkspaceID   string
	stacksCreateGcpIntegrationID string
	stacksCreateGCSBucket        string
	stacksCreateGCSKey           string
)

var stacksCmd = &cobra.Command{
	Use:   "stacks [cmd]",
	Short: "Manage Stacks",
	Args:  cobra.MinimumNArgs(1),
}

var stacksCreateCmd = &cobra.Command{
	Use:           "create ENVIRONMENT_ID NAME",
	Short:         "Create a Stack",
	Args:          cobra.ExactArgs(2),
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, args []string) error {
		stack, err := c.CreateStack(args[0], args[1])
		if err != nil {
			return fmt.Errorf("failed creating stack: %w", err)
		}

		fmt.Fprintf(os.Stdout, "Created stack %s\n", stack.ID)

		var policy []byte

		if stacksCreateStatePath != "" {
			policy, err = os.ReadFile(stacksCreateStatePath)
			if err != nil {
				return fmt.Errorf("failed reading state file: %w", err)
			}
		} else if stacksCreateLocation == "s3" && stacksCreateS3IntegrationID != "" {
			policy, err = json.Marshal(map[string]interface{}{
				"awsIntegration": stacksCreateS3IntegrationID,
				"s3Bucket":       stacksCreateS3Bucket,
				"s3Key":          stacksCreateS3Key,
			})
			if err != nil {
				return fmt.Errorf("failed encoding S3 policy: %w", err)
			}
		} else if stacksCreateLocation == "tfc" && stacksCreateTfCAPIToken != "" {
			policy, err = json.Marshal(map[string]interface{}{
				"apiToken":    stacksCreateTfCAPIToken,
				"workspaceId": stacksCreateTfCWorkspaceID,
			})
			if err != nil {
				return fmt.Errorf("failed encoding TfC policy: %w", err)
			}
		} else if stacksCreateLocation == "gcs" && stacksCreateGcpIntegrationID != "" {
			policy, err = json.Marshal(map[string]interface{}{
				"gcpIntegration": stacksCreateGcpIntegrationID,
				"gcsBucket":      stacksCreateGCSBucket,
				"gcsKey":         stacksCreateGCSKey,
			})
			if err != nil {
				return fmt.Errorf("failed encoding Gcs policy: %w", err)
			}
		}

		if len(policy) > 0 {
			err = c.UpdateStatePolicy(
				stack.ID,
				stacksCreateLocation,
				json.RawMessage(policy),
			)
			if err != nil {
				return fmt.Errorf("failed updating state policy: %w", err)
			}

			fmt.Fprintf(os.Stdout, "Updated state policy %s\n", stacksCreateStatePath)
		}

		return nil
	},
}

var stacksListCmd = &cobra.Command{
	Use:           "list <env_id>",
	Short:         "List States in an Environment",
	Args:          cobra.ExactArgs(1),
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, args []string) error {
		list, err := c.ListStacks(args[0])
		if err != nil {
			return fmt.Errorf("failed listing stacks: %w", err)
		}

		return render(list)
	},
}

var stacksGetCmd = &cobra.Command{
	Use:           "get <environment_id> <stack_id>",
	Short:         "Get a Specific Stack",
	Args:          cobra.ExactArgs(2),
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, args []string) error {
		stack, err := c.GetStack(args[0], args[1])
		if err != nil {
			return fmt.Errorf("failed getting stack %s:%s: %w", args[0], args[1], err)
		}

		return render(stack)
	},
}

var stacksDeleteCmd = &cobra.Command{
	Use:           "delete <environment_id> <stack_id>",
	Short:         "Delete a Specific Stack",
	Args:          cobra.ExactArgs(2),
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, args []string) error {
		stack, err := c.DeleteStack(args[0], args[1])
		if err != nil {
			return fmt.Errorf("failed deleting stack %s:%s: %w", args[0], args[1], err)
		}

		return render(stack)
	},
}

// nolint: lll
func init() {
	stacksCreateCmd.PersistentFlags().StringVar(&stacksCreateStatePath, "state-path", "", "Path to state file")
	stacksCreateCmd.PersistentFlags().StringVar(&stacksCreateLocation, "location", "manual", "Policy location")
	stacksCreateCmd.PersistentFlags().StringVar(&stacksCreateS3IntegrationID, "s3-aws-integration-id", "", "AWS Integration ID")
	stacksCreateCmd.PersistentFlags().StringVar(&stacksCreateS3Bucket, "s3-bucket", "", "S3 bucket ARN")
	stacksCreateCmd.PersistentFlags().StringVar(&stacksCreateS3Key, "s3-key", "", "S3 key")
	stacksCreateCmd.PersistentFlags().StringVar(&stacksCreateTfCAPIToken, "tfc-api-token", "", "Terraform Cloud API Token")
	stacksCreateCmd.PersistentFlags().StringVar(&stacksCreateTfCWorkspaceID, "tfc-workspace-id", "", "Terraform Cloud workspace ID")
	stacksCreateCmd.PersistentFlags().StringVar(&stacksCreateGcpIntegrationID, "gcp-integration-id", "", "GCP Integration ID")
	stacksCreateCmd.PersistentFlags().StringVar(&stacksCreateGCSBucket, "gcs-bucket", "", "GCS name")
	stacksCreateCmd.PersistentFlags().StringVar(&stacksCreateGCSKey, "gcs-key", "", "GCS key")
	stacksCmd.AddCommand(stacksCreateCmd, stacksListCmd, stacksGetCmd, stacksDeleteCmd)
	rootCmd.AddCommand(stacksCmd)
}
