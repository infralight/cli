package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
)

var (
	stacksCreateName             string
	stacksCreateEnv              string
	stacksCreateLocation         string
	stacksCreateStatePath        string
	stacksCreateS3IntegrationID  string
	stacksCreateS3Bucket         string
	stacksCreateTfCAPIToken      string
	stacksCreateTfCWorkspaceID   string
	stacksCreateGcpIntegrationID string
	stacksCreateGCSBucket        string
	stacksListEnv                string
	stacksGetEnv                 string
	StacksGetStackId             string
	stacksDeleteEnv              string
	StacksDeleteStackId          string
)

var VALID_LOCATIONS = []string{"s3", "gcs", "tfc", "local"}

var stacksCmd = &cobra.Command{
	Use:   "stacks [cmd]",
	Short: "Manage Stacks",
	Args:  cobra.MinimumNArgs(1),
}

var stacksCreateCmd = &cobra.Command{
	Use:           "create --env-id ENVIRONMENT_ID --name NAME --state-path STATE_PATH --location FILE_LOCATION",
	Short:         "Create a Stack",
	Args:          cobra.NoArgs,
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, args []string) error {
		sort.Strings(VALID_LOCATIONS)
		if ind := sort.SearchStrings(VALID_LOCATIONS, stacksCreateLocation); ind > len(VALID_LOCATIONS) || ind < 0 || stacksCreateLocation != VALID_LOCATIONS[ind] {
			return fmt.Errorf("invalid location - %s", stacksCreateLocation)
		}
		stack, err := c.CreateStack(stacksCreateEnv, stacksCreateName)
		if err != nil {
			return fmt.Errorf("failed creating stack: %w", err)
		}

		fmt.Fprintf(os.Stdout, "Created stack %s\n", stack.ID)

		var policy []byte

		if stacksCreateLocation == "s3" && stacksCreateS3IntegrationID != "" {
			policy, err = json.Marshal(map[string]interface{}{
				"awsIntegration": stacksCreateS3IntegrationID,
				"s3Bucket":       stacksCreateS3Bucket,
				"s3Key":          stacksCreateStatePath,
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
				"gcsKey":         stacksCreateStatePath,
			})
			if err != nil {
				return fmt.Errorf("failed encoding Gcs policy: %w", err)
			}
		} else if stacksCreateLocation == "local" {
			policy, err = os.ReadFile(stacksCreateStatePath)
			if err != nil {
				return fmt.Errorf("failed reading state file: %w", err)
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
	Use:           "list --env-id ENVIRONMENT_ID",
	Short:         "List States in an Environment",
	Args:          cobra.NoArgs,
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, args []string) error {
		list, err := c.ListStacks(stacksListEnv)
		if err != nil {
			return fmt.Errorf("failed listing stacks: %w", err)
		}

		return render(list)
	},
}

var stacksGetCmd = &cobra.Command{
	Use:           "get --env-id ENVIRONMENT_ID --stack-id STACK_ID",
	Short:         "Get a Specific Stack",
	Args:          cobra.NoArgs,
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, args []string) error {
		stack, err := c.GetStack(stacksGetEnv, StacksGetStackId)
		if err != nil {
			return fmt.Errorf("failed getting stack %s:%s: %w", args[0], args[1], err)
		}

		return render(stack)
	},
}

var stacksDeleteCmd = &cobra.Command{
	Use:           "delete --env-id ENVIRONMENT_ID --stack-id STACK_ID",
	Short:         "Delete a Specific Stack",
	Args:          cobra.NoArgs,
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, args []string) error {
		stack, err := c.DeleteStack(stacksDeleteEnv, StacksDeleteStackId)
		if err != nil {
			return fmt.Errorf("failed deleting stack %s:%s: %w", args[0], args[1], err)
		}

		return render(stack)
	},
}

// nolint
func init() {
	stacksCreateCmd.PersistentFlags().StringVar(&stacksCreateName, "name", "", "Stack name")
	stacksCreateCmd.MarkPersistentFlagRequired("name")
	stacksCreateCmd.PersistentFlags().StringVar(&stacksCreateEnv, "env-id", "", "The environment id to create the stack into")
	stacksCreateCmd.MarkPersistentFlagRequired("env-id")
	stacksCreateCmd.PersistentFlags().StringVar(&stacksCreateStatePath, "state-path", "", "Path to state file (whether local or remote location")
	stacksCreateCmd.MarkPersistentFlagRequired("state-path")
	stacksCreateCmd.PersistentFlags().StringVar(&stacksCreateLocation, "location", "local", "Policy location [s3, gcp, tfc, local]")
	stacksCreateCmd.MarkPersistentFlagRequired("location")
	stacksCreateCmd.PersistentFlags().StringVar(&stacksCreateS3IntegrationID, "s3-aws-integration-id", "", "AWS Integration ID")
	stacksCreateCmd.PersistentFlags().StringVar(&stacksCreateS3Bucket, "s3-bucket", "", "S3 bucket ARN")
	stacksCreateCmd.PersistentFlags().StringVar(&stacksCreateStatePath, "s3-key", "", "S3 key")
	stacksCreateCmd.PersistentFlags().StringVar(&stacksCreateTfCAPIToken, "tfc-api-token", "", "Terraform Cloud API Token")
	stacksCreateCmd.PersistentFlags().StringVar(&stacksCreateTfCWorkspaceID, "tfc-workspace-id", "", "Terraform Cloud workspace ID")
	stacksCreateCmd.PersistentFlags().StringVar(&stacksCreateGcpIntegrationID, "gcp-integration-id", "", "GCP Integration ID")
	stacksCreateCmd.PersistentFlags().StringVar(&stacksCreateGCSBucket, "gcs-bucket", "", "GCS name")
	stacksCreateCmd.PersistentFlags().StringVar(&stacksCreateStatePath, "gcs-key", "", "GCS key")

	stacksGetCmd.PersistentFlags().StringVar(&StacksGetStackId, "stack-id", "", "Stack Id")
	stacksGetCmd.MarkPersistentFlagRequired("stack-id")
	stacksGetCmd.PersistentFlags().StringVar(&stacksGetEnv, "env-id", "", "Environment Id of the chosen stack")
	stacksGetCmd.MarkPersistentFlagRequired("env-id")

	stacksDeleteCmd.PersistentFlags().StringVar(&StacksDeleteStackId, "stack-id", "", "Stack Id")
	stacksDeleteCmd.MarkPersistentFlagRequired("stack-id")
	stacksDeleteCmd.PersistentFlags().StringVar(&stacksDeleteEnv, "env-id", "", "Environment Id of the chosen stack")
	stacksDeleteCmd.MarkPersistentFlagRequired("env-id")

	stacksListCmd.PersistentFlags().StringVar(&stacksListEnv, "env-id", "", "The environment id to list the stacks for")
	stacksListCmd.MarkPersistentFlagRequired("env-id")
	stacksCmd.AddCommand(stacksCreateCmd, stacksListCmd, stacksGetCmd, stacksDeleteCmd)
	rootCmd.AddCommand(stacksCmd)
}
