package steamstore

import (
	"context"

	"github.com/google/go-querystring/query"
)

func (c *Client) GetSteamSpyAppDetails(ctx context.Context, appId uint) (*SteamSpyAppDetailsResponse, error) {
	var res *SteamSpyAppDetailsResponse

	q, _ := query.Values(&SteamSpyQuery{
		Request: "appdetails",
		AppId:   appId,
	})
	err := c.get(ctx, SteamSpyApiBaseUrl+"?"+q.Encode(), &res, false)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) GetSteamSpyAppsPaginated(ctx context.Context, page uint) (map[string]SteamSpyAppDetailsResponse, error) {
	var res map[string]SteamSpyAppDetailsResponse

	q, _ := query.Values(SteamSpyQuery{
		Request: "all",
		Page:    page,
	})
	if err := c.get(ctx, SteamSpyApiBaseUrl+"?"+q.Encode(), &res, false); err != nil {
		return nil, err
	}

	return res, nil
}
