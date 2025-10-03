package scraper

import (
	"fmt"
	"io"
	"net/http"
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

func (c *Client) Get(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.userAgent)

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
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			bytes, err := io.ReadAll(resp.Body)

			if err != nil {
				return nil, fmt.Errorf("failed to read response body: %w", err)
			}

			return bytes, nil
		}

		if i < c.retries {
			time.Sleep(time.Duration(i+1) * time.Second)
		}
	}
	if lastErr != nil {
		return nil, fmt.Errorf("failed to fetch url %s after %d retries: - last error: %w", url, c.retries, lastErr)
	}
	return nil, fmt.Errorf("failed to fetch url %s after %d retries: final status %d", url, c.retries, lastStatusCode)
}
