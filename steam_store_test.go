package steamstore

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

const envSteamApiKey = "STEAM_API_KEY"

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

func getKeyFromEnv(t *testing.T) string {
	key, ok := os.LookupEnv(envSteamApiKey)
	if !ok {
		t.Fatalf("getKeyFromEnv() %s key not found", envSteamApiKey)
	}
	return key
}

func Test_GetAppList(t *testing.T) {
	c := New(
		WithKey(getKeyFromEnv(t)),
	)

	q := AppListQuery{
		MaxResults: 1,
	}

	res, err := c.GetAppList(context.TODO(), q)
	if err != nil {
		t.Fatalf("GetAppList() error = %v", err)
	}

	if len(res.Response.Apps) == 0 {
		t.Fatalf("GetAppList() response apps count = %v, want %v", len(res.Response.Apps), q.MaxResults)
	}

	if res.Response.Apps[0].AppID != 10 {
		t.Errorf("GetAppList() AppID = %v, want %v", res.Response.Apps[0].AppID, 10)
	}
}

func Test_GetTagList(t *testing.T) {
	c := New(
		WithKey(getKeyFromEnv(t)),
	)

	q := TagListQuery{
		Language: "english",
	}

	res, err := c.GetTagList(context.TODO(), q)
	if err != nil {
		t.Fatalf("GetTagList() error = %v", err)
	}

	want := map[uint]string{
		9:  "Strategy",
		19: "Action",
		21: "Adventure",
	}
	for _, tag := range res.Response.Tags {
		if _, ok := want[tag.TagId]; ok {
			delete(want, tag.TagId)
		}

		if len(want) == 0 {
			break
		}
	}

	if len(want) != 0 {
		t.Errorf("GetTagList() missing tags = %v", want)
	}
}

func Test_GetAppDetails(t *testing.T) {
	c := New()

	res, err := c.GetAppDetails(context.TODO(), 10)
	if err != nil {
		t.Fatalf("GetAppDetails() error = %v", err)
	}

	app, ok := res["10"]
	if !ok {
		t.Fatalf("GetAppDetails() response app = %v, want app id %v", res, 10)
	}

	if app.Data.SteamAppId != 10 {
		t.Errorf("GetAppDetails() SteamAppId = %v, want %v", app.Data.SteamAppId, 10)
	}
}

func Test_GetSteamSpyAppDetails(t *testing.T) {
	c := New()

	res, err := c.GetSteamSpyAppDetails(context.TODO(), 10)
	if err != nil {
		t.Fatalf("GetSteamSpyAppDetails() error = %v", err)
	}

	if res.AppId != 10 {
		t.Errorf("GetSteamSpyAppDetails() AppId = %v, want %v", res.AppId, 10)
	}
}

func Test_GetSteamSpyAppDetails_emptyTags(t *testing.T) {
	c := New()

	jsonFile, err := os.Open("test_files/steam_spy_app_details_empty_tags.json")
	if err != nil {
		t.Fatalf("error opening JSON test file: %v", err)
	}
	defer jsonFile.Close()

	mockData, err := io.ReadAll(jsonFile)
	if err != nil {
		t.Fatalf("error reading JSON test file: %v", err)
	}

	rs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(mockData)
	}))
	defer rs.Close()

	var res *SteamSpyAppDetailsResponse
	err = c.get(context.Background(), rs.URL, &res, false)
	if err != nil {
		t.Fatalf("GetSteamSpyAppDetails_emptyTags() error = %v", err)
	}

	if res.AppId != 1620 {
		t.Errorf("GetSteamSpyAppDetails_emptyTags() AppId = %v, want %v", res.AppId, 1620)
	}
}
