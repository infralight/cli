package client

import (
	"fmt"

	"github.com/ido50/requests"
)

const (
	DefaultInfralightURL = "https://prodapi.infralight.cloud/api"
	DefaultAuthHeader    = "Authorization"
)

type Client struct {
	authHeader string
	httpc      *requests.HTTPClient
}

func New(url, authHeader string) *Client {
	if url == "" {
		url = DefaultInfralightURL
	}
	if authHeader == "" {
		authHeader = DefaultAuthHeader
	}

	return &Client{
		authHeader: authHeader,
		httpc:      requests.NewClient(url).Accept("application/json"),
	}
}

func (c *Client) Authenticate(accessKey, secretKey string) error {
	var creds struct {
		Token     string `json:"access_token"`
		ExpiresIn int64  `json:"expires_in"`
		Type      string `json:"token_type"`
	}

	err := c.httpc.NewRequest("POST", "/account/access_keys/login").
		JSONBody(map[string]interface{}{
			"accessKey": accessKey,
			"secretKey": secretKey,
		}).
		Into(&creds).
		Run()
	if err != nil {
		return err
	}

	c.httpc.Header(c.authHeader, fmt.Sprintf("Bearer %s", creds.Token))

	return nil
}

type Stack struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (c *Client) ListStacks() (list []Stack, err error) {
	err = c.httpc.NewRequest("GET", "/stacks").
		Into(&list).
		Run()
	return list, err
}
