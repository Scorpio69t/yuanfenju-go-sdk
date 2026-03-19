package yuanfenju

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const defaultBaseURL = "https://api.yuanfenju.com/index.php"

type Config struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
	Timeout    time.Duration
}

type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client

	Free       *FreeService
	Bazi       *BaziService
	Divination *DivinationService
}

func NewClient(cfg Config) (*Client, error) {
	if strings.TrimSpace(cfg.APIKey) == "" {
		return nil, errors.New("api key is required")
	}

	baseURL := strings.TrimSpace(cfg.BaseURL)
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	baseURL = strings.TrimRight(baseURL, "/")

	hc := cfg.HTTPClient
	if hc == nil {
		timeout := cfg.Timeout
		if timeout <= 0 {
			timeout = 10 * time.Second
		}
		hc = &http.Client{Timeout: timeout}
	}

	c := &Client{
		apiKey:     cfg.APIKey,
		baseURL:    baseURL,
		httpClient: hc,
	}

	c.Free = &FreeService{client: c}
	c.Bazi = &BaziService{client: c}
	c.Divination = &DivinationService{client: c}

	return c, nil
}

type CommonResponse[T any] struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	Notice  string `json:"notice"`
	Data    T      `json:"data"`
}

type APIError struct {
	Code    int
	Message string
	Notice  string
}

func (e *APIError) Error() string {
	if e == nil {
		return ""
	}
	if e.Notice == "" {
		return fmt.Sprintf("yuanfenju api error: code=%d message=%s", e.Code, e.Message)
	}
	return fmt.Sprintf("yuanfenju api error: code=%d message=%s notice=%s", e.Code, e.Message, e.Notice)
}

func (c *Client) doForm(ctx context.Context, path string, form url.Values, out any) error {
	if form == nil {
		form = url.Values{}
	}
	if form.Get("api_key") == "" {
		form.Set("api_key", c.apiKey)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4*1024))
		return fmt.Errorf("http %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return err
	}

	return nil
}
