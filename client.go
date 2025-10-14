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
)

var (
	ErrNoApiKey = errors.New("api key is required")
)

type Client struct {
	Opts
}

type OptFunc func(opts *Opts)

type Opts struct {
	key     string
	timeout time.Duration
}

func defaultOpts() Opts {
	return Opts{
		timeout: time.Second * 10,
	}
}

func WithKey(key string) OptFunc {
	return func(opts *Opts) {
		opts.key = key
	}
}

func WithTimeout(timeout time.Duration) OptFunc {
	return func(opts *Opts) {
		opts.timeout = timeout
	}
}

func NewClient(opts ...OptFunc) *Client {
	o := defaultOpts()
	for _, opt := range opts {
		opt(&o)
	}

	return &Client{o}
}

func (c *Client) SetKey(key string) {
	c.key = key
}

func (c *Client) get(ctx context.Context, url string, output any, key bool) error {
	if key && c.key == "" {
		return ErrNoApiKey
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	res, err := client.Do(req)
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
