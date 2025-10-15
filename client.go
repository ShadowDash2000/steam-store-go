package steamstore

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

const (
	SteamApiBaseUrl       = "https://api.steampowered.com"
	SteamShadowApiBaseUrl = "https://store.steampowered.com/api"
	SteamSpyApiBaseUrl    = "https://steamspy.com/api.php"

	defaultTimeout = time.Second * 10
	defaultRate    = 4
)

var (
	ErrNoApiKey = errors.New("api key is required")
)

type Client struct {
	Opts
}

type OptFunc func(opts *Opts)

type Opts struct {
	client  *http.Client
	limiter *rateLimiter
	key     string
}

func defaultOpts() Opts {
	return Opts{
		client: &http.Client{
			Timeout: defaultTimeout,
		},
		limiter: newRateLimiter(defaultRate),
	}
}

func WithKey(key string) OptFunc {
	return func(opts *Opts) {
		opts.key = key
	}
}

func WithLimiter(rate int) OptFunc {
	return func(opts *Opts) {
		opts.limiter = newRateLimiter(rate)
	}
}

func WithTimeout(timeout int) OptFunc {
	return func(opts *Opts) {
		if timeout > 0 {
			opts.client.Timeout = time.Duration(timeout) * time.Second
		} else {
			opts.client.Timeout = defaultTimeout
		}
	}
}

func New(opts ...OptFunc) *Client {
	c := &Client{
		Opts: defaultOpts(),
	}

	for _, opt := range opts {
		opt(&c.Opts)
	}

	return c
}

func (c *Client) SetKey(key string) {
	c.key = key
}

func (c *Client) get(ctx context.Context, url string, output any, key bool) error {
	if key && c.key == "" {
		return ErrNoApiKey
	}

	c.limiter.wait()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		return errors.New(string(body))
	}

	err = json.NewDecoder(res.Body).Decode(&output)
	if err != nil {
		return err
	}

	return nil
}
