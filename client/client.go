package client

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
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
		httpc:      requests.NewClient(url).Accept("application/json"),
	}
}

func (c *Client) Authenticate(accessKey, secretKey string) error {
	// check if we have an authentication token already stored in the
	// temporary directory. Files must only be valid for 3 minutes, after
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

	if time.Since(stat.ModTime()) > 3*time.Minute {
		// file too old, remove it
		return os.Remove(path)
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

func (c *Client) Codify(assetType, assetID string) (output interface{}, err error) {
	err = c.httpc.NewRequest("POST", "/reverseLearning").
		JSONBody(map[string]string{
			"assetType": assetType,
			"assetId":   assetID,
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
