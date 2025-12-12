package steamstore

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/go-querystring/query"
)

func (c *Client) GetSteamSpyAppDetails(ctx context.Context, appId uint) (*SteamSpyAppDetailsResponse, error) {
	var res *SteamSpyAppDetailsResponse

	q, _ := query.Values(&SteamSpyQuery{
		Request: "appdetails",
		AppId:   appId,
	})
	_, err := c.get(ctx, SteamSpyApiBaseUrl+"?"+q.Encode(), &res, false)
	if err != nil {
		return nil, err
	}

	return res, nil
}

var ErrSteamSpyLastPage = errors.New("steam-store: no more pages")

func (c *Client) GetSteamSpyAppsPaginated(ctx context.Context, page uint) (map[string]SteamSpyAppDetailsResponse, error) {
	var res map[string]SteamSpyAppDetailsResponse

	q, _ := query.Values(SteamSpyQuery{
		Request: "all",
		Page:    page,
	})
	if statusCode, err := c.get(ctx, SteamSpyApiBaseUrl+"?"+q.Encode(), &res, false); err != nil {
		if statusCode == http.StatusInternalServerError {
			return nil, ErrSteamSpyLastPage
		}
		return nil, err
	}

	return res, nil
}
