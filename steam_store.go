package steamstore

import (
	"context"
	"fmt"

	"github.com/google/go-querystring/query"
)

func (c *Client) GetAppList(ctx context.Context, opts AppListQuery) (*AppListResponse, error) {
	var appListRes *AppListResponse

	q, _ := query.Values(opts)
	err := c.get(ctx, SteamApiBaseUrl+"/IStoreService/GetAppList/v2/?"+q.Encode(), &appListRes, true)
	if err != nil {
		return nil, err
	}

	return appListRes, nil
}

func (c *Client) GetTagList(ctx context.Context, opts TagListQuery) (*TagListResponse, error) {
	var tagListRes *TagListResponse

	q, _ := query.Values(opts)
	err := c.get(ctx, SteamApiBaseUrl+"/IStoreService/GetTagList/v1/"+q.Encode(), &tagListRes, true)
	if err != nil {
		return nil, err
	}

	return tagListRes, nil
}

func (c *Client) GetAppDetails(ctx context.Context, appId uint) (*AppDetailsResponse, error) {
	var appDetailsRes *AppDetailsResponse

	err := c.get(ctx, fmt.Sprintf(SteamShadowApiBaseUrl+"?appids=%d", appId), &appDetailsRes, false)
	if err != nil {
		return nil, err
	}

	return appDetailsRes, nil
}

func (c *Client) GetSteamSpyAppDetails(ctx context.Context, appId uint) (*SteamSpyAppDetailsResponse, error) {
	var appDetailsRes *SteamSpyAppDetailsResponse

	q, _ := query.Values(&SteamSpyQuery{
		Request: "appdetails",
		AppId:   appId,
	})
	err := c.get(ctx, SteamSpyApiBaseUrl+"?"+q.Encode(), &appDetailsRes, false)
	if err != nil {
		return nil, err
	}

	return appDetailsRes, nil
}

func (c *Client) GetAllApps(ctx context.Context) chan App {
	ch := make(chan App)

	go func() {
		query := AppListQuery{
			IncludeGames: true,
			MaxResults:   10000,
		}

		for {
			select {
			case <-ctx.Done():
				close(ch)
				return
			default:
				res, err := c.GetAppList(ctx, query)
				if err != nil {
					continue
				}

				for _, app := range res.Response.Apps {
					ch <- app
				}

				if !res.Response.HaveMoreResults {
					close(ch)
					return
				}

				query.LastAppId = res.Response.LastAppId
			}
		}
	}()

	return ch
}
