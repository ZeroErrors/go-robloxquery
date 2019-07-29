package list

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

type Request struct {
	SortToken                  *string
	GameFilter                 *string
	TimeFilter                 *string
	GenreFilter                *string
	ExclusiveStartId           *int64
	SortOrder                  *int
	GameSetTargetId            *int64
	Keyword                    *string
	StartRows                  *int
	MaxRows                    *int
	IsKeywordSuggestionEnabled *bool
	ContextCountryRegionId     *int
	ContextUniverseId          *int64
}

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

func Do(request Request) (Response, error) {
	var response Response

	valueOfRequest := reflect.ValueOf(request)
	requestType := valueOfRequest.Type()

	var query bytes.Buffer
	for i := 0; i < valueOfRequest.NumField(); i++ {
		field := valueOfRequest.Field(i)
		if field.IsNil() {
			continue
		}
		field = field.Elem()
		query.WriteString("model.")
		query.WriteString(requestType.Field(i).Name)
		query.WriteString("=\"")
		query.WriteString(strings.ReplaceAll(fmt.Sprintf("%v", field.Interface()), "\"", "\\\""))
		query.WriteString("\"")
	}

	// Note: the '_' field is required for the request to go work for some reason
	req, err := http.NewRequest("GET", "https://games.roblox.com/v1/games/list?_=-1&"+query.String(), nil)
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
