package steamstore

import (
	"context"
	"errors"
	"net/http"
	"time"

	"golang.org/x/time/rate"
	"resty.dev/v3"
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
	client  *resty.Client
	limiter *rate.Limiter
	key     string
}

func defaultOpts() Opts {
	limiter := rate.NewLimiter(rate.Every(2*time.Second), defaultRate)
	cb := resty.NewCircuitBreaker()

	return Opts{
		client: resty.New().
			SetCircuitBreaker(cb).
			SetTimeout(defaultTimeout).
			AddRequestMiddleware(func(client *resty.Client, request *resty.Request) error {
				if err := limiter.Wait(request.Context()); err != nil {
					return err
				}
				return nil
			}),
		limiter: limiter,
	}
}

func WithKey(key string) OptFunc {
	return func(opts *Opts) {
		opts.key = key
	}
}

func WithRateLimit(rate rate.Limit) OptFunc {
	return func(opts *Opts) {
		opts.limiter.SetLimit(rate)
	}
}

func WithTimeout(timeout int) OptFunc {
	return func(opts *Opts) {
		if timeout > 0 {
			opts.client.SetTimeout(time.Duration(timeout) * time.Second)
		}
	}
}

func WithRetryCount(count int) OptFunc {
	return func(opts *Opts) {
		opts.client.SetRetryCount(count)
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

func (c *Client) get(ctx context.Context, url string, output any, needKey bool) error {
	if needKey && c.key == "" {
		return ErrNoApiKey
	}

	req := c.client.R().
		SetHeader("Accept", "application/json").
		SetContext(ctx).
		SetResult(output)

	if needKey {
		req.SetQueryParam("key", c.key)
	}

	res, err := req.Execute(http.MethodGet, url)
	if err != nil {
		return err
	}

	if res.IsError() {
		return errors.New("steam-store.get(): " + res.String())
	}

	return nil
}
