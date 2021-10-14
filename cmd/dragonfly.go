package cmd

import (
	"fmt"

	"github.com/infralight/cli/client"
	"github.com/infralight/cli/version"
	"github.com/spf13/cobra"
)

var dflyCmd = &cobra.Command{
	Use:   "pipelines [cmd]",
	Short: "Manage Dragonfly Pipelines",
	Args:  cobra.MinimumNArgs(1),
}

var pipelineListCmd = &cobra.Command{
	Use:           "list",
	Short:         "List Pipelines",
	Args:          cobra.NoArgs,
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, _ []string) error {
		list, err := c.Dragonfly.ListPipelines()
		if err != nil {
			return fmt.Errorf("failed listing pipelines: %w", err)
		}

		return render(list)
	},
}

var (
	pipelineCreateInput client.CreatePipelineInput
	pipelineCreateCmd   = &cobra.Command{
		Use:           "create <params>",
		Short:         "Create a Pipeline",
		Args:          cobra.NoArgs,
		SilenceErrors: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			env, err := c.Dragonfly.CreatePipeline(pipelineCreateInput)
			if err != nil {
				return fmt.Errorf("failed creating pipeline: %w", err)
			}

			return render(env)
		},
	}
)

var (
	pipelineGetCmd = &cobra.Command{
		Use:           "get PIPELINE_ID",
		Short:         "Get a Pipeline",
		Args:          cobra.ExactArgs(1),
		SilenceErrors: true,
		RunE: func(_ *cobra.Command, args []string) error {
			env, err := c.Dragonfly.GetPipeline(args[0])
			if err != nil {
				return fmt.Errorf("failed fetching pipeline: %w", err)
			}

			return render(env)
		},
	}
)

var (
	pipelineUpdateInput client.UpdatePipelineInput
	pipelineUpdateCmd   = &cobra.Command{
		Use:           "update PIPELINE_ID",
		Short:         "Update a Pipeline",
		Args:          cobra.ExactArgs(1),
		SilenceErrors: true,
		RunE: func(_ *cobra.Command, args []string) error {
			env, err := c.Dragonfly.UpdatePipeline(args[0], pipelineUpdateInput)
			if err != nil {
				return fmt.Errorf("failed updating pipeline: %w", err)
			}

			return render(env)
		},
	}
)

var (
	pipelineDeleteCmd = &cobra.Command{
		Use:           "delete PIPELINE_ID",
		Short:         "Delete a Pipeline",
		Args:          cobra.ExactArgs(1),
		SilenceErrors: true,
		RunE: func(_ *cobra.Command, args []string) error {
			env, err := c.Dragonfly.DeletePipeline(args[0])
			if err != nil {
				return fmt.Errorf("failed deleting pipeline: %w", err)
			}

			return render(env)
		},
	}
)

var (
	pipelineLockCmd = &cobra.Command{
		Use:           "lock PIPELINE_ID",
		Short:         "Lock a Pipeline",
		Args:          cobra.ExactArgs(1),
		SilenceErrors: true,
		RunE: func(_ *cobra.Command, args []string) error {
			env, err := c.Dragonfly.LockPipeline(args[0])
			if err != nil {
				return fmt.Errorf("failed locking pipeline: %w", err)
			}

			return render(env)
		},
	}
)

var (
	pipelineUnlockCmd = &cobra.Command{
		Use:           "unlock PIPELINE_ID",
		Short:         "Unlock a Pipeline",
		Args:          cobra.ExactArgs(1),
		SilenceErrors: true,
		RunE: func(_ *cobra.Command, args []string) error {
			env, err := c.Dragonfly.UnlockPipeline(args[0])
			if err != nil {
				return fmt.Errorf("failed unlocking pipeline: %w", err)
			}

			return render(env)
		},
	}
)

var (
	envVarAddInput client.AddPipelineVariableInput
	envVarAddCmd   = &cobra.Command{
		Use:           "add-variable PIPELINE_ID",
		Short:         "Add an Environment Variable",
		Args:          cobra.ExactArgs(1),
		SilenceErrors: true,
		RunE: func(_ *cobra.Command, args []string) error {
			env, err := c.Dragonfly.AddPipelineVariable(args[0], envVarAddInput)
			if err != nil {
				return fmt.Errorf("failed adding environment variable: %w", err)
			}

			return render(env)
		},
	}
)

var envVarListCmd = &cobra.Command{
	Use:           "list-variables PIPELINE_ID",
	Short:         "List Pipelines Environment Variables",
	Args:          cobra.ExactArgs(1),
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, args []string) error {
		list, err := c.Dragonfly.ListPipelineVariables(args[0])
		if err != nil {
			return fmt.Errorf("failed listing pipeline variables: %w", err)
		}

		return render(list)
	},
}

var (
	envVarUpdateInput client.UpdatePipelineVariableInput
	envVarUpdateCmd   = &cobra.Command{
		Use:           "update-variable PIPELINE_ID VARIABLE_ID",
		Short:         "Update a Pipeline's Environment Variable",
		Args:          cobra.ExactArgs(2),
		SilenceErrors: true,
		RunE: func(_ *cobra.Command, args []string) error {
			err := c.Dragonfly.UpdatePipelineVariable(
				args[0],
				args[1],
				envVarUpdateInput,
			)
			if err != nil {
				return fmt.Errorf("failed updating pipeline variable: %w", err)
			}

			return nil
		},
	}
)

var envVarDeleteCmd = &cobra.Command{
	Use:           "delete-variable PIPELINE_ID VARIABLE_ID",
	Short:         "Delete a Pipeline's Environment Variable",
	Args:          cobra.ExactArgs(2),
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, args []string) error {
		variable, err := c.Dragonfly.DeletePipelineVariable(args[0], args[1])
		if err != nil {
			return fmt.Errorf("failed deleting pipeline variable: %w", err)
		}

		return render(variable)
	},
}

var policyListCmd = &cobra.Command{
	Use:           "list-policies PIPELINE_ID",
	Short:         "List Pipeline Policies",
	Args:          cobra.ExactArgs(1),
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, args []string) error {
		list, err := c.Dragonfly.ListPipelinePolicies(args[0])
		if err != nil {
			return fmt.Errorf("failed listing pipeline policies: %w", err)
		}

		return render(list)
	},
}

var (
	policyAddInput client.AddPipelinePolicyInput
	policyAddCmd   = &cobra.Command{
		Use:           "add-policy PIPELINE_ID",
		Short:         "Add a Pipeline Policy",
		Args:          cobra.ExactArgs(1),
		SilenceErrors: true,
		RunE: func(_ *cobra.Command, args []string) error {
			policy, err := c.Dragonfly.AddPipelinePolicy(args[0], policyAddInput)
			if err != nil {
				return fmt.Errorf("failed adding policy: %w", err)
			}

			return render(policy)
		},
	}
)

var (
	policyDeleteCmd = &cobra.Command{
		Use:           "delete-policy PIPELINE_ID POLICY_ID",
		Short:         "Delete a Pipeline Policy",
		Args:          cobra.ExactArgs(2),
		SilenceErrors: true,
		RunE: func(_ *cobra.Command, args []string) error {
			policy, err := c.Dragonfly.DeletePipelinePolicy(args[0], args[1])
			if err != nil {
				return fmt.Errorf("failed deleting policy: %w", err)
			}

			return render(policy)
		},
	}
)

var (
	runsListInput client.ListPipelineRunsInput
	runsListCmd   = &cobra.Command{
		Use:           "list-runs PIPELINE_ID",
		Short:         "List Pipeline Executions",
		Args:          cobra.ExactArgs(1),
		SilenceErrors: true,
		RunE: func(_ *cobra.Command, args []string) error {
			list, err := c.Dragonfly.ListPipelineRuns(args[0], runsListInput)
			if err != nil {
				return fmt.Errorf("failed listing pipeline executions: %w", err)
			}

			return render(list)
		},
	}
)

var runsGetCmd = &cobra.Command{
	Use:           "get-run PIPELINE_ID RUN_ID",
	Short:         "Get Pipeline Execution",
	Args:          cobra.ExactArgs(2),
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, args []string) error {
		run, err := c.Dragonfly.GetPipelineRun(args[0], args[1])
		if err != nil {
			return fmt.Errorf("failed getting pipeline execution: %w", err)
		}

		return render(run)
	},
}

var (
	runsCreateInput client.CreatePipelineRunInput
	runsCreateCmd   = &cobra.Command{
		Use:           "run PIPELINE_ID [params]",
		Short:         "Execute a Pipeline",
		Args:          cobra.ExactArgs(1),
		SilenceErrors: true,
		RunE: func(_ *cobra.Command, args []string) error {
			env, err := c.Dragonfly.CreatePipelineRun(args[0], runsCreateInput)
			if err != nil {
				return fmt.Errorf("failed executing pipeline: %w", err)
			}

			return render(env)
		},
	}
)

func init() { // nolint: funlen
	// Pipeline creation flags
	pipelineCreateCmd.PersistentFlags().StringVar(
		&pipelineCreateInput.Name, "name", "", "Pipeline name (required)",
	)
	pipelineCreateCmd.MarkPersistentFlagRequired("name") // nolint: errcheck
	pipelineCreateCmd.PersistentFlags().StringVar(
		&pipelineCreateInput.Owner,
		"owner",
		fmt.Sprintf("%sCli", version.Product),
		fmt.Sprintf("Pipeline owner. Use `%s users list` to see available options", version.Product),
	)
	pipelineCreateCmd.PersistentFlags().StringSliceVar(
		&pipelineCreateInput.Labels, "labels", nil, "Pipeline labels",
	)

	// Pipeline modification flags
	pipelineUpdateCmd.PersistentFlags().StringVar(
		&pipelineUpdateInput.Name, "name", "", "Pipeline name (required)",
	)
	pipelineUpdateCmd.MarkPersistentFlagRequired("name") // nolint: errcheck
	pipelineUpdateCmd.PersistentFlags().StringVar(
		&pipelineUpdateInput.Owner,
		"owner",
		fmt.Sprintf("%sCli", version.Product),
		fmt.Sprintf("Pipeline owner. Use `%s users list` to see available options", version.Product),
	)
	pipelineUpdateCmd.PersistentFlags().BoolVar(
		&pipelineUpdateInput.IsTerragrunt,
		"is-terragrunt",
		false,
		"Is Terragrunt? (simply provide flag to indicate true or use --is-terragrunt=true|false)",
	)
	pipelineUpdateCmd.PersistentFlags().BoolVar(
		&pipelineUpdateInput.AutoApprove,
		"auto-approve",
		false,
		"Auto Approve (simply provide flag to indicate true or use --auto-approve=true|false)",
	)
	pipelineUpdateCmd.PersistentFlags().StringVar(
		&pipelineUpdateInput.TerraformVersion,
		"terraform-version",
		"",
		"Terrafrom Version (required)",
	)
	pipelineUpdateCmd.MarkPersistentFlagRequired("terraform-version") // nolint: errcheck
	pipelineUpdateCmd.PersistentFlags().StringVar(
		&pipelineUpdateInput.VCS.Branch,
		"vcs-branch",
		"master",
		"VCS Branch",
	)
	pipelineUpdateCmd.PersistentFlags().StringVar(
		&pipelineUpdateInput.VCS.Path,
		"vcs-path",
		"/",
		"VCS Path",
	)
	pipelineUpdateCmd.PersistentFlags().StringVar(
		&pipelineUpdateInput.VCS.Repo,
		"vcs-repo",
		"",
		"VCS Repository (required)",
	)
	pipelineUpdateCmd.MarkPersistentFlagRequired("vcs-repo") // nolint: errcheck
	pipelineUpdateCmd.PersistentFlags().StringVar(
		&pipelineUpdateInput.VCS.VCSID,
		"vcs-id",
		"",
		"VCS ID (required)",
	)
	pipelineUpdateCmd.MarkPersistentFlagRequired("vcs-id") // nolint: errcheck
	pipelineUpdateCmd.PersistentFlags().BoolVar(
		&pipelineUpdateInput.FailurePolicy.FailOnWarn,
		"policy-fail-on-warn",
		false,
		"Fail on Warn? (simply provide flag to indicate true or use --policy-fail-on-warn=true|false)",
	)
	pipelineUpdateCmd.PersistentFlags().BoolVar(
		&pipelineUpdateInput.FailurePolicy.NoFail,
		"policy-no-fail",
		false,
		"No Fail (simply provide flag to indicate true or use --policy-no-fail=true|false)",
	)

	// Environment Variable creation flags
	envVarAddCmd.PersistentFlags().StringVar(
		&envVarAddInput.Type,
		"type",
		"var",
		"Variable Type",
	)
	envVarAddCmd.PersistentFlags().StringVar(
		&envVarAddInput.Key,
		"key",
		"",
		"Variable Key (required)",
	)
	envVarAddCmd.MarkPersistentFlagRequired("key") // nolint: errcheck
	envVarAddCmd.PersistentFlags().StringVar(
		&envVarAddInput.Value,
		"value",
		"",
		"Variable Value (required)",
	)
	envVarAddCmd.MarkPersistentFlagRequired("value") // nolint: errcheck
	envVarAddCmd.PersistentFlags().StringVar(
		&envVarAddInput.Description,
		"description",
		"",
		"Variable Description",
	)
	envVarAddCmd.PersistentFlags().BoolVar(
		&envVarAddInput.Encrypted,
		"encrypted",
		false,
		"Is the Value Encrypted? (simply provide flag to indicate true or use --encrypted=true|false)",
	)

	// Environment variable update flags
	envVarUpdateCmd.PersistentFlags().StringVar(
		&envVarUpdateInput.Key,
		"key",
		"",
		"Variable Key (required)",
	)
	envVarUpdateCmd.MarkPersistentFlagRequired("key") // nolint: errcheck
	envVarUpdateCmd.PersistentFlags().StringVar(
		&envVarUpdateInput.Value,
		"value",
		"",
		"Variable Value (required)",
	)
	envVarUpdateCmd.MarkPersistentFlagRequired("value") // nolint: errcheck
	envVarUpdateCmd.PersistentFlags().StringVar(
		&envVarUpdateInput.Description,
		"description",
		"",
		"Variable Description",
	)

	// Add policy parameters
	policyAddCmd.PersistentFlags().StringVar(
		&policyAddInput.RuleID,
		"rule-id",
		"",
		"Rule ID (required)",
	)
	policyAddCmd.MarkPersistentFlagRequired("rule-id") // nolint: errcheck
	policyAddCmd.PersistentFlags().BoolVar(
		&policyAddInput.IsPostPlan,
		"is-post-plan",
		false,
		"Is Post Plan? (simply provide flag to indicate true or use --is-post-plan=true|false)",
	)

	// List runs parameters
	runsListCmd.PersistentFlags().IntVar(
		&runsListInput.Page,
		"page",
		0,
		"Page",
	)
	runsListCmd.PersistentFlags().IntVar(
		&runsListInput.Limit,
		"limit",
		0,
		"Limit",
	)

	// Create run parameters
	runsCreateCmd.PersistentFlags().StringVar(
		&runsCreateInput.UserID,
		"user-id",
		"",
		"User ID (required)",
	)
	runsCreateCmd.MarkPersistentFlagRequired("user-id") // nolint: errcheck
	runsCreateCmd.PersistentFlags().StringVar(
		&runsCreateInput.Description,
		"description",
		"Triggered by CLI",
		"Execution Description",
	)

	dflyCmd.AddCommand(
		pipelineListCmd,
		pipelineCreateCmd,
		pipelineGetCmd,
		pipelineUpdateCmd,
		pipelineDeleteCmd,
		pipelineLockCmd,
		pipelineUnlockCmd,
		envVarListCmd,
		envVarAddCmd,
		envVarUpdateCmd,
		envVarDeleteCmd,
		policyListCmd,
		policyAddCmd,
		policyDeleteCmd,
		runsListCmd,
		runsGetCmd,
		runsCreateCmd,
	)
	rootCmd.AddCommand(dflyCmd)
}
