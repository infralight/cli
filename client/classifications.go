package client

import (
	"encoding/base64"
	"fmt"
	"time"
)

type ClassificationsClient struct {
	*baseClient
}

type Classification struct {
	ID          string    `json:"id,omitempty"`
	ModifiedID  string    `json:"_id,omitempty"`
	AccountID   string    `json:"-"`
	Type        string    `json:"type"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Labels      []string  `json:"labels,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`

	// The classification's Rego policy
	Rego string `json:"rego"`

	// If true, the Rego value is base64 encoded
	IsRegoEncoded bool `json:"isRegoEncoded"`
}

type ListClassificationsInput struct {
	DecodeRego bool
}

func (c *ClassificationsClient) List(input ListClassificationsInput) (
	list []Classification,
	err error,
) {
	err = c.httpc.NewRequest("GET", "/classifications").
		Into(&list).
		Run()
	if err != nil {
		return list, err
	}

	if input.DecodeRego {
		for i := range list {
			if list[i].Rego != "" {
				decoded, err := base64.StdEncoding.DecodeString(list[i].Rego)
				if err != nil {
					// ignore the error, we'll just return the policy still base64
					// encoded
					list[i].IsRegoEncoded = true
				} else {
					list[i].Rego = string(decoded)
					list[i].IsRegoEncoded = false
				}
			}
		}
	}

	return list, nil
}

type CreateClassificationInput struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Rego        string   `json:"code"`
	Labels      []string `json:"labels"`
}

func (c *ClassificationsClient) Create(input CreateClassificationInput) (
	class Classification,
	err error,
) {
	err = c.httpc.NewRequest("POST", "/classifications").
		JSONBody(input).
		Into(&class).
		Run()
	return class, err
}

type UpdateClassificationInput struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Rego        string   `json:"code"`
	Labels      []string `json:"labels"`
}

func (c *ClassificationsClient) Update(id string, input UpdateClassificationInput) (
	class Classification,
	err error,
) {
	err = c.httpc.NewRequest("PUT", fmt.Sprintf("/classifications/%s", id)).
		JSONBody(input).
		Into(&class).
		Run()
	return class, err
}

func (c *ClassificationsClient) Delete(id string) (err error) {
	err = c.httpc.
		NewRequest("DELETE", fmt.Sprintf("/classifications/%s", id)).
		Run()
	return err
}

type ClassificationResult struct {
	ID   string `json:"id"`
	ARN  string `json:"arn"`
	Type string `json:"type"`
}

func (c *ClassificationsClient) Run(id string) (
	res []ClassificationResult,
	err error,
) {
	// we need to first get the classification. Unfortunately, the app-server
	// doesn't currently have an API to get a specific classification, so we
	// must run the List method and extract it from the response
	list, err := c.List(ListClassificationsInput{
		DecodeRego: false,
	})
	if err != nil {
		return res, fmt.Errorf("failed listing rules: %w", err)
	}

	var rule Classification
	for i := range list {
		if list[i].ID == id {
			rule = list[i]
			break
		}
	}

	if rule.ID != id {
		return res, fmt.Errorf("no such rule %q", id)
	}

	err = c.httpc.
		NewRequest("POST", "/exclusions/test").
		QueryParam("version", "2").
		JSONBody(map[string]interface{}{
			"name":        rule.Name,
			"description": rule.Description,
			"type":        rule.Type,
			"code":        rule.Rego,
		}).
		Into(&res).
		Run()
	if err != nil {
		return res, fmt.Errorf("failed running rule: %w", err)
	}

	return res, nil
}
