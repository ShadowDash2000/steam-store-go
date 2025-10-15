# Steam Store GO
Steam Store GO uses: official, unofficial Steam API and SteamSpy API.

#### Install:

```bash
go get -u github.com/ShadowDash2000/steam-store-go
```

#### How to use:

```go
package main

import (
	"context"

	steamstore "github.com/ShadowDash2000/steam-store-go"
)

func main() {
	steam := steamstore.New(
		steamstore.WithKey("PROVIDE_STEAM_API_KEY_HERE"),
	)

	res, err := steam.GetAppDetails(context.Background(), 10)
	if err != nil {
		// error handling
	}

	// ...
}
```

### Steam API Key

Proceed to this link to obtain API Key: https://steamcommunity.com/dev/apikey
