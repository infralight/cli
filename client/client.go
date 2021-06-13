package client

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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
		httpc: requests.NewClient(url).
			Accept("application/json").
			ErrorHandler(func(status int, _ string, body io.Reader) error {
				var errMap map[string]string
				err := json.NewDecoder(body).Decode(&errMap)
				if err == nil {
					if msg, ok := errMap["message"]; ok {
						return errors.New(msg)
					}
				}

				return fmt.Errorf("server returned unexpected status %d", status)
			}),
	}
}

func (c *Client) Authenticate(accessKey, secretKey string) error {
	// check if we have an authentication token already stored in the
	// temporary directory. Files must only be valid for 5 minutes, after
	// which they should be removed and recreated.
	err := c.getCachedToken(accessKey)
	if err == nil {
		return nil
	}

	var creds struct {
		Token     string `json:"access_token"`
		ExpiresIn int64  `json:"expires_in"`
		Type      string `json:"token_type"`
	}

	err = c.httpc.NewRequest("POST", "/account/access_keys/login").
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

	return c.saveCachedToken(accessKey, creds.Token)
}

func (c *Client) cacheName(accessKey string) string {
	return filepath.Join(os.TempDir(), fmt.Sprintf("%x", sha256.Sum224([]byte(accessKey))))
}

func (c *Client) getCachedToken(accessKey string) error {
	path := c.cacheName(accessKey)

	stat, err := os.Stat(path)
	if err != nil {
		return err
	}

	if time.Since(stat.ModTime()) > 5*time.Minute {
		// file too old, remove it
		os.Remove(path) // nolint: errcheck
		return errors.New("cached token too old")
	}

	token, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	token = bytes.TrimSpace(token)

	c.httpc.Header(c.authHeader, fmt.Sprintf("Bearer %s", token))
	return nil
}

func (c *Client) saveCachedToken(accessKey, token string) error {
	return os.WriteFile(c.cacheName(accessKey), []byte(token), 0600)
}

type Environment struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Owner string `json:"owner,omitempty"`
}

func (c *Client) ListEnvironments() (list []Environment, err error) {
	err = c.httpc.NewRequest("GET", "/environments").
		Into(&list).
		Run()
	return list, err
}

type Stack struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (c *Client) ListStacks(envID string) (list []Stack, err error) {
	err = c.httpc.NewRequest("GET", fmt.Sprintf("/environments/%s/stack/", envID)).
		Into(&list).
		Run()
	return list, err
}

func (c *Client) GetStack(envID, stackID string) (stack map[string]interface{}, err error) {
	err = c.httpc.NewRequest("GET", fmt.Sprintf("/environments/%s/stack/%s", envID, stackID)).
		Into(&stack).
		Run()
	return stack, err
}

func (c *Client) DeleteStack(envID, stackID string) (stack map[string]interface{}, err error) {
	err = c.httpc.NewRequest("DELETE", fmt.Sprintf("/environments/%s/stack/%s", envID, stackID)).
		Into(&stack).
		Run()
	return stack, err
}

func (c *Client) Codify(assetType, assetID string) (output string, err error) {
	err = c.httpc.NewRequest("POST", "/reverseLearning").
		JSONBody(map[string]string{
			"assetType": assetType,
			"assetId":   assetID,
		}).
		BodyHandler(func(_ int, contentType string, body io.Reader, target interface{}) error {
			b, err := io.ReadAll(body)
			if err != nil {
				return err
			}

			if len(b) == 0 {
				return errors.New("no content received from server")
			}

			v := target.(*string)
			*v, err = strconv.Unquote(string(b))
			if err != nil {
				*v = string(b)
				return nil
			}

			*v, err = strconv.Unquote(*v)
			if err != nil {
				*v = string(b)
			}

			return nil
		}).
		Into(&output).
		Run()
	return output, err
}

type Drift struct {
	ID        string `json:"driftId"`
	CreatedAt int64  `json:"createdAt"`
	Total     int64  `json:"total"`
	Unmanaged int64  `json:"unmanagedCount"`
	Managed   int64  `json:"managedCount"`
	Modified  int64  `json:"modifiedCount"`
}

func (d Drift) CreationDate() time.Time {
	return time.Unix(d.CreatedAt, 0)
}

type AssetState string

const (
	StateManaged   AssetState = "managed"
	StateModified  AssetState = "modified"
	StateUnmanaged AssetState = "unmanaged"
)

type Asset struct {
	ID            string          `json:"assetId"`
	Type          string          `json:"assetType"`
	Hash          string          `json:"hash"`
	State         AssetState      `json:"state"`
	InventoryItem json.RawMessage `json:"inventoryItem,omitempty"`
	StateItem     json.RawMessage `json:"stateItem,omitempty"`
}

func (c *Client) ListDrifts(onlyDelta bool, limit uint64) (list []Drift, err error) {
	req := c.httpc.NewRequest("GET", "/drifts").
		QueryParam("onlyDelta", strconv.FormatBool(onlyDelta)).
		Into(&list)

	if limit > 0 {
		req.QueryParam("limit", fmt.Sprintf("%d", limit))
	}

	err = req.Run()
	return list, err
}

func (c *Client) ShowDrift(driftID string) (list []Asset, err error) {
	err = c.httpc.NewRequest(
		"GET",
		fmt.Sprintf("/drifts/%s", strings.TrimPrefix(driftID, "Drifts/")),
	).
		Into(&list).
		Run()
	return list, err
}

func (c *Client) ShowAsset(assetID string) (list []Asset, err error) {
	err = c.httpc.NewRequest(
		"GET",
		fmt.Sprintf("/drifts/asset/%s", url.PathEscape(assetID)),
	).
		Into(&list).
		Run()
	return list, err
}

type State struct {
	ID        string          `json:"id"`
	StackID   string          `json:"stackId"`
	CreatedAt time.Time       `json:"createdAt,omitempty"`
	Policy    json.RawMessage `json:"policy"`
	RunID     string          `json:"runId,omitempty"`
}

func (c *Client) ListStates(stackID string) (list []State, err error) {
	err = c.httpc.NewRequest(
		"GET",
		fmt.Sprintf("/states/stack/%s", url.PathEscape(stackID)),
	).
		Into(&list).
		Run()
	return list, err
}

func (c *Client) GetLatestState(stackID string) (list State, err error) {
	err = c.httpc.NewRequest(
		"GET",
		fmt.Sprintf("/states/stack/%s/latest", url.PathEscape(stackID)),
	).
		Into(&list).
		Run()
	return list, err
}

func (c *Client) UploadStatePolicy(
	stackID string,
	tfState, policy json.RawMessage,
) (err error) {
	var jsonPolicy map[string]interface{}
	err = json.Unmarshal(policy, &jsonPolicy)
	if err != nil {
		return fmt.Errorf("failed decoding policy: %w", err)
	}

	err = c.httpc.NewRequest(
		"POST",
		fmt.Sprintf("/states/stack/%s/upload", url.PathEscape(stackID)),
	).
		JSONBody(map[string]interface{}{
			"tfState": string(tfState),
			"policy":  jsonPolicy,
		}).
		ExpectedStatus(http.StatusNoContent).
		Run()
	return err
}
