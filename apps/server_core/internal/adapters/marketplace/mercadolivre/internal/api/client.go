// Package api is Mercado Livre's wire: HTTP, auth header, pagination tokens
// and raw DTOs. Nothing outside adapters/marketplace/mercadolivre can import
// it — that is the Go internal rule at the vendor root (§2.2), and it is the
// boundary the legacy connectors module never had.
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// TokenSource supplies a live access token. The adapter does not know where
// tokens come from (env, DB credential, future account context) — the
// composition root decides.
type TokenSource func(ctx context.Context) (string, error)

type Client struct {
	base   string
	userID string
	token  TokenSource
	http   *http.Client
}

func NewClient(baseURL, userID string, token TokenSource) (*Client, error) {
	if baseURL == "" || userID == "" || token == nil {
		return nil, fmt.Errorf("mercadolivre api: base url, user id and token source are all required")
	}
	return &Client{base: baseURL, userID: userID, token: token, http: &http.Client{Timeout: 30 * time.Second}}, nil
}

func (c *Client) getJSON(ctx context.Context, path string, out any) error {
	tok, err := c.token(ctx)
	if err != nil {
		return fmt.Errorf("mercadolivre api: token: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.base+path, nil)
	if err != nil {
		return fmt.Errorf("mercadolivre api: build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+tok)
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("mercadolivre api: GET %s: %w", path, err)
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("mercadolivre api: read %s: %w", path, err)
	}
	if resp.StatusCode != http.StatusOK {
		// O status e o corpo dizem o quê; o token NUNCA aparece em erro.
		return fmt.Errorf("mercadolivre api: GET %s: status %d: %s", path, resp.StatusCode, truncate(body, 300))
	}
	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("mercadolivre api: decode %s: %w", path, err)
	}
	return nil
}

func truncate(b []byte, n int) string {
	if len(b) <= n {
		return string(b)
	}
	return string(b[:n]) + "…"
}
