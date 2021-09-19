package client

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type DragonflyClient struct {
	*baseClient
}

type Pipeline struct {
	ID               string                `json:"id"`
	AccountID        string                `json:"-"`
	Name             string                `json:"name"`
	AutoApprove      bool                  `json:"autoApprove"`
	Configured       bool                  `json:"configured"`
	IsLocked         bool                  `json:"isLocked"`
	IsTerragrunt     bool                  `json:"isTerragrunt"`
	Labels           []string              `json:"labels"`
	Owner            string                `json:"owner"`
	RunID            string                `json:"runId"`
	TerraformVersion string                `json:"terraformVersion"`
	CreatedAt        time.Time             `json:"createdAt"`
	UpdatedAt        time.Time             `json:"updatedAt"`
	VCS              PipelineVCS           `json:"vcs,omitempty"`
	FailurePolicy    PipelineFailurePolicy `json:"policy,omitempty"`
}

type PipelineVCS struct {
	Branch string `json:"branch"`
	Path   string `json:"path"`
	Repo   string `json:"repo"`
	VCSID  string `json:"vcsId"`
}

type PipelineFailurePolicy struct {
	FailOnWarn bool `json:"failOnWarn"`
	NoFail     bool `json:"noFail"`
}

type PipelinePolicy struct {
	DeletedID  string    `json:"_id,omitempty"`
	ID         string    `json:"id,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	RuleID     string    `json:"ruleId"`
	IsPostPlan bool      `json:"isPostPlan"`
}

type PipelineRun struct {
	ID            string    `json:"id"`
	AccountID     string    `json:"-"`
	EnvironmentID string    `json:"environmentId"`
	UserID        string    `json:"userId"`
	Source        string    `json:"source"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	Status        string    `json:"status"`
	Branch        string    `json:"branch"`
	Commit        string    `json:"commit"`
	Description   string    `json:"description"`
	HasChanges    bool      `json:"hasChanges"`
}

func (c *DragonflyClient) ListPipelines() (pipelines []Pipeline, err error) {
	err = c.httpc.NewRequest("GET", "/dragonfly/environments").
		Into(&pipelines).
		Run()
	return pipelines, err
}

type CreatePipelineInput struct {
	Name   string   `json:"name"`
	Owner  string   `json:"owner"`
	Labels []string `json:"labels"`
}

func (c *DragonflyClient) CreatePipeline(
	input CreatePipelineInput,
) (pipeline Pipeline, err error) {
	err = c.httpc.NewRequest("POST", "/dragonfly/environments").
		JSONBody(input).
		Into(&pipeline).
		Run()
	return pipeline, err
}

func (c *DragonflyClient) GetPipeline(pipelineID string) (
	pipeline Pipeline,
	err error,
) {
	err = c.httpc.
		NewRequest(
			"GET",
			fmt.Sprintf("/dragonfly/environments/%s", pipelineID),
		).
		Into(&pipeline).
		Run()
	return pipeline, err
}

type UpdatePipelineInput struct {
	Name             string                `json:"name"`
	Owner            string                `json:"owner"`
	IsTerragrunt     bool                  `json:"isTerragrunt"`
	AutoApprove      bool                  `json:"autoApprove"`
	TerraformVersion string                `json:"terraformVersion"`
	VCS              PipelineVCS           `json:"vcs,omitempty"`
	FailurePolicy    PipelineFailurePolicy `json:"policy,omitempty"`
}

func (c *DragonflyClient) UpdatePipeline(
	pipelineID string,
	input UpdatePipelineInput,
) (pipeline Pipeline, err error) {
	err = c.httpc.
		NewRequest(
			"PUT",
			fmt.Sprintf("/dragonfly/environments/%s", pipelineID),
		).
		JSONBody(input).
		Into(&pipeline).
		Run()
	return pipeline, err
}

func (c *DragonflyClient) DeletePipeline(pipelineID string) (
	deletedPipeline Pipeline,
	err error,
) {
	err = c.httpc.
		NewRequest(
			"DELETE",
			fmt.Sprintf("/dragonfly/environments/%s", pipelineID),
		).
		Into(&deletedPipeline).
		Run()
	return deletedPipeline, err
}

func (c *DragonflyClient) LockPipeline(envID string) (
	env Pipeline,
	err error,
) {
	err = c.httpc.
		NewRequest(
			"POST",
			fmt.Sprintf("/dragonfly/environments/%s/lock", envID),
		).
		JSONBody(map[string]interface{}{"lock": true}).
		Into(&env).
		Run()
	return env, err
}

func (c *DragonflyClient) UnlockPipeline(envID string) (
	env Pipeline,
	err error,
) {
	err = c.httpc.
		NewRequest(
			"POST",
			fmt.Sprintf("/dragonfly/environments/%s/lock", envID),
		).
		JSONBody(map[string]interface{}{"lock": false}).
		Into(&env).
		Run()
	return env, err
}

type PipelineVariable struct {
	DeletedID   string    `json:"_id,omitempty"`
	ID          string    `json:"id,omitempty"`
	AccountID   string    `json:"-"`
	PipelineID  string    `json:"environmentId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Type        string    `json:"type"`
	Description string    `json:"description,omitempty"`
	Encrypted   bool      `json:"encrypted"`
	Key         string    `json:"key"`
	Value       string    `json:"value"`
}

func (c *DragonflyClient) ListPipelineVariables(pipelineID string) (
	variables []PipelineVariable,
	err error,
) {
	err = c.httpc.
		NewRequest(
			"GET",
			fmt.Sprintf("/dragonfly/environments/%s/variables", pipelineID),
		).
		Into(&variables).
		Run()
	return variables, err
}

type AddPipelineVariableInput struct {
	Type        string `json:"type"`
	Key         string `json:"key"`
	Value       string `json:"value"`
	Encrypted   bool   `json:"encrypted"`
	Description string `json:"description"`
}

type AddPipelineVariableOutput struct {
	ID string `json:"_id"`
}

func (c *DragonflyClient) AddPipelineVariable(
	pipelineID string,
	input AddPipelineVariableInput,
) (output AddPipelineVariableOutput, err error) {
	err = c.httpc.
		NewRequest(
			"POST",
			fmt.Sprintf("/dragonfly/environments/%s/variables", pipelineID),
		).
		JSONBody(input).
		Into(&output).
		Run()
	return output, err
}

type UpdatePipelineVariableInput struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

func (c *DragonflyClient) UpdatePipelineVariable(
	pipelineID string,
	variableID string,
	input UpdatePipelineVariableInput,
) (err error) {
	err = c.httpc.
		NewRequest(
			"PUT",
			fmt.Sprintf(
				"/dragonfly/environments/%s/variables/%s",
				pipelineID,
				variableID,
			),
		).
		JSONBody(input).
		ExpectedStatus(http.StatusNoContent).
		Run()
	return err
}

func (c *DragonflyClient) DeletePipelineVariable(pipelineID, varID string) (
	output PipelineVariable,
	err error,
) {
	err = c.httpc.
		NewRequest(
			"DELETE",
			fmt.Sprintf("/dragonfly/environments/%s/variables/%s", pipelineID, varID),
		).
		Into(&output).
		Run()
	return output, err
}

func (c *DragonflyClient) ListPipelinePolicies(pipelineID string) (policies []PipelinePolicy, err error) {
	err = c.httpc.
		NewRequest(
			"GET",
			fmt.Sprintf("/dragonfly/environments/%s/policies", pipelineID),
		).
		Into(&policies).
		Run()
	return policies, err
}

type AddPipelinePolicyInput struct {
	RuleID     string `json:"ruleId"`
	IsPostPlan bool   `json:"isPostPlan"`
}

func (c *DragonflyClient) AddPipelinePolicy(
	pipelineID string,
	input AddPipelinePolicyInput,
) (output interface{}, err error) {
	err = c.httpc.
		NewRequest(
			"POST",
			fmt.Sprintf("/dragonfly/environments/%s/policies", pipelineID),
		).
		JSONBody(input).
		Into(&output).
		Run()
	return output, err
}

func (c *DragonflyClient) DeletePipelinePolicy(pipelineID, policyID string) (
	output PipelinePolicy,
	err error,
) {
	err = c.httpc.
		NewRequest(
			"DELETE",
			fmt.Sprintf("/dragonfly/environments/%s/policies/%s", pipelineID, policyID),
		).
		Into(&output).
		Run()
	return output, err
}

type ListPipelineRunsInput struct {
	Page  int
	Limit int
}

func (c *DragonflyClient) ListPipelineRuns(
	pipelineID string,
	input ListPipelineRunsInput,
) (runs []PipelineRun, err error) {
	req := c.httpc.
		NewRequest(
			"GET",
			fmt.Sprintf("/dragonfly/environments/%s/runs", pipelineID),
		).
		Into(&runs)

	if input.Page > 0 {
		req.QueryParam("page", strconv.Itoa(input.Page))
	}
	if input.Limit > 0 {
		req.QueryParam("limit", strconv.Itoa(input.Limit))
	}

	err = req.Run()
	return runs, err
}

type CreatePipelineRunInput struct {
	UserID      string `json:"userId"`
	Description string `json:"description"`
}

func (c *DragonflyClient) CreatePipelineRun(
	pipelineID string,
	input CreatePipelineRunInput,
) (run PipelineRun, err error) {
	err = c.httpc.
		NewRequest(
			"POST",
			fmt.Sprintf("/dragonfly/environments/%s/runs", pipelineID),
		).
		JSONBody(input).
		Into(&run).
		Run()
	return run, err
}

func (c *DragonflyClient) GetPipelineRun(pipelineID, runID string) (
	run PipelineRun,
	err error,
) {
	err = c.httpc.
		NewRequest(
			"GET",
			fmt.Sprintf("/dragonfly/environments/%s/runs/%s", pipelineID, runID),
		).
		Into(&run).
		Run()
	return run, err
}
