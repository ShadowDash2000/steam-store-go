package steamstore

import (
	"context"
	"fmt"
	"strconv"

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

func (c *Client) GetSteamSpyAppDetails(ctx context.Context, appId uint) (*SteamSpyAppDetailsResponse, error) {
	var resRaw *SteamSpyAppDetailsResponseRaw

	q, _ := query.Values(&SteamSpyQuery{
		Request: "appdetails",
		AppId:   appId,
	})
	err := c.get(ctx, SteamSpyApiBaseUrl+"?"+q.Encode(), &resRaw, false)
	if err != nil {
		return nil, err
	}

	res, err := c.formatSteamSpyAppDetailsResponse(resRaw)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) formatSteamSpyAppDetailsResponse(resRaw *SteamSpyAppDetailsResponseRaw) (*SteamSpyAppDetailsResponse, error) {
	var err error
	res := &SteamSpyAppDetailsResponse{
		AppId:          resRaw.AppId,
		Name:           resRaw.Name,
		Developer:      resRaw.Developer,
		Publisher:      resRaw.Publisher,
		ScoreRank:      resRaw.ScoreRank,
		Positive:       resRaw.Positive,
		Negative:       resRaw.Negative,
		UserScore:      resRaw.UserScore,
		Owners:         resRaw.Owners,
		AverageForever: resRaw.AverageForever,
		Average2Weeks:  resRaw.Average2Weeks,
		MedianForever:  resRaw.MedianForever,
		Median2Weeks:   resRaw.Median2Weeks,
		CCU:            resRaw.CCU,
		Languages:      resRaw.Languages,
		Genre:          resRaw.Genre,
		Tags:           resRaw.Tags,
	}

	res.Price, err = strconv.ParseUint(resRaw.Price, 10, 32)
	if err != nil {
		return nil, err
	}
	res.InitialPrice, err = strconv.ParseUint(resRaw.InitialPrice, 10, 32)
	if err != nil {
		return nil, err
	}
	res.Discount, err = strconv.ParseUint(resRaw.Discount, 10, 32)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) GetAllApps(ctx context.Context) chan App {
	ch := make(chan App)

	go func() {
		q := AppListQuery{
			IncludeGames: true,
			MaxResults:   10000,
		}

		for {
			select {
			case <-ctx.Done():
				close(ch)
				return
			default:
				res, err := c.GetAppList(ctx, q)
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

				q.LastAppId = res.Response.LastAppId
			}
		}
	}()

	return ch
}
