package steamstore

import (
	"context"
	"fmt"

	"github.com/google/go-querystring/query"
)

func (c *Client) GetAppList(ctx context.Context, opts AppListQuery) (*AppListResponse, error) {
	var res *AppListResponse

	q, _ := query.Values(opts)
	q.Set("key", c.key)
	err := c.get(ctx, SteamApiBaseUrl+"/IStoreService/GetAppList/v1/?"+q.Encode(), &res, true)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) GetTagList(ctx context.Context, opts TagListQuery) (*TagListResponse, error) {
	var res *TagListResponse

	q, _ := query.Values(opts)
	q.Set("key", c.key)
	err := c.get(ctx, SteamApiBaseUrl+"/IStoreService/GetTagList/v1/?"+q.Encode(), &res, true)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) GetAppDetails(ctx context.Context, appId uint) (AppDetailsResponse, error) {
	var res AppDetailsResponse

	err := c.get(ctx, fmt.Sprintf(SteamShadowApiBaseUrl+"/appdetails?appids=%d", appId), &res, false)
	if err != nil {
		return nil, err
	}

	return res, nil
}

type GetAllAppsMessage struct {
	Apps []App
	Err  error
}

func (c *Client) GetAllApps(ctx context.Context, opts AppListQuery) chan GetAllAppsMessage {
	ch := make(chan GetAllAppsMessage, 1)

	var lastAppId uint
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(ch)
				return
			default:
				if lastAppId > 0 {
					opts.LastAppId = lastAppId
				}

				res, err := c.GetAppList(ctx, opts)
				if err != nil {
					ch <- GetAllAppsMessage{Err: err}
					continue
				}

				ch <- GetAllAppsMessage{Apps: res.Response.Apps}

				if !res.Response.HaveMoreResults {
					close(ch)
					return
				}

				lastAppId = res.Response.LastAppId
			}
		}
	}()

	return ch
}
