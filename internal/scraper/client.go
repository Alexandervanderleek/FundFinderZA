package scraper

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	httpClient *http.Client
	retries    int
	userAgent  string
}

type ClientOption func(*Client)

func WithRetries(retries int) ClientOption {
	return func(c *Client) {
		c.retries = retries
	}
}

func WithUserAgent(userAgent string) ClientOption {
	return func(c *Client) {
		c.userAgent = userAgent
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

const defaultUserAgent = "placeholder"
const defaultRetries = 3
const defaultTimeout = 30 * time.Second

func NewClient(options ...ClientOption) *Client {
	client := &Client{
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		retries:   defaultRetries,
		userAgent: defaultUserAgent,
	}

	for _, option := range options {
		option(client)
	}

	return client
}

func (c *Client) Post(url string, formData url.Values) ([]byte, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(formData.Encode()))

	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", url)

	return c.doRequest(req)
}

func (c *Client) Get(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.userAgent)

	return c.doRequest(req)
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	var lastErr error
	var lastStatusCode int

	for i := 0; i < c.retries; i++ {
		resp, err := c.httpClient.Do(req)

		if err != nil {
			lastErr = err
			if i < c.retries {
				time.Sleep(time.Duration(i+1) * time.Second)
				continue
			}
			break
		}
		defer resp.Body.Close()

		lastStatusCode = resp.StatusCode

		if resp.StatusCode == http.StatusOK {
			bytes, err := io.ReadAll(resp.Body)

			if err != nil {
				return nil, fmt.Errorf("failed to read response body: %w", err)
			}

			return bytes, nil
		}

		if i < c.retries-1 {
			time.Sleep(time.Duration(i+1) * time.Second)
		}
	}

	if lastErr != nil {
		return nil, fmt.Errorf("failed to fetch after %d retries: - last error: %w", c.retries, lastErr)
	}
	return nil, fmt.Errorf("failed to fetch after %d retries: final status %d", c.retries, lastStatusCode)
}
