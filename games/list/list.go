package list

import (
	"encoding/json"
	"github.com/google/go-querystring/query"
	"io/ioutil"
	"net/http"
)

// Request specifies the data to request
type Request struct {
	Timestamp                  int64   `url:"_"`
	SortToken                  *string `url:"model.sortToken,omitempty"`
	GameFilter                 *string `url:"model.gameFilter,omitempty"`
	TimeFilter                 *string `url:"model.timeFilter,omitempty"`
	GenreFilter                *string `url:"model.genreFilter,omitempty"`
	ExclusiveStartID           *int64  `url:"model.exclusiveStartId,omitempty"`
	SortOrder                  *int    `url:"model.sortOrder,omitempty"`
	GameSetTargetID            *int64  `url:"model.gameSetTargetId,omitempty"`
	Keyword                    *string `url:"model.keyword,omitempty"`
	StartRows                  *int    `url:"model.startRows,omitempty"`
	MaxRows                    *int    `url:"model.maxRows,omitempty"`
	IsKeywordSuggestionEnabled *bool   `url:"model.isKeywordSuggestionEnabled,omitempty"`
	ContextCountryRegionID     *int    `url:"model.contextCountryRegionId,omitempty"`
	ContextUniverseID          *int64  `url:"model.contextUniverseId,omitempty"`
}

// Response specifies the data that will be received as the reply
type Response struct {
	Games []struct {
		CreatorID      int    `json:"creatorId"`
		CreatorName    string `json:"creatorName"`
		CreatorType    string `json:"creatorType"`
		TotalUpVotes   int    `json:"totalUpVotes"`
		TotalDownVotes int    `json:"totalDownVotes"`
		UniverseID     int    `json:"universeId"`
		Name           string `json:"name"`
		PlaceID        int    `json:"placeId"`
		PlayerCount    int    `json:"playerCount"`
		ImageToken     string `json:"imageToken"`
		Users          []struct {
			UserID int    `json:"userId"`
			GameID string `json:"gameId"`
		} `json:"users"`
		IsSponsored         bool    `json:"isSponsored"`
		NativeAdData        *string `json:"nativeAdData"`
		Price               *int    `json:"price"`
		AnalyticsIdentifier *string `json:"analyticsIdentifier"`
	} `json:"games"`
	SuggestedKeyword         string `json:"suggestedKeyword"`
	CorrectedKeyword         string `json:"correctedKeyword"`
	FilteredKeyword          string `json:"filteredKeyword"`
	HasMoreRows              bool   `json:"hasMoreRows"`
	NextPageExclusiveStartID int    `json:"nextPageExclusiveStartId"`
}

// Do makes a request with the provided Request data and returns either a Response or an error
func Do(request Request) (Response, error) {
	var response Response

	v, err := query.Values(request)
	if err != nil {
		return response, err
	}

	req, err := http.NewRequest("GET", "https://games.roblox.com/v1/games/list?"+v.Encode(), nil)
	if err != nil {
		return response, err
	}

	req.Header.Add("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return response, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(body, &response)
	return response, err
}
