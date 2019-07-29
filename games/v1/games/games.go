package list

import (
	"encoding/json"
	"github.com/google/go-querystring/query"
	"io/ioutil"
	"net/http"
	"time"
)

// Request specifies the data to request
type Request struct {
	Timestamp   int64   `url:"_"`
	UniverseIDs []int64 `url:"universeIds,comma,omitempty"`
}

// Response specifies the data that will be received as the reply
type Response struct {
	Data []struct {
		ID          int    `json:"id"`
		RootPlaceID int    `json:"rootPlaceId"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Creator     struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"creator"`
		Price                     int       `json:"price"`
		IsExperimental            bool      `json:"isExperimental"`
		AllowedGearGenres         []string  `json:"allowedGearGenres"`
		AllowedGearCategories     []string  `json:"allowedGearCategories"`
		Playing                   int       `json:"playing"`
		Visits                    int       `json:"visits"`
		MaxPlayers                int       `json:"maxPlayers"`
		Created                   time.Time `json:"created"`
		Updated                   time.Time `json:"updated"`
		StudioAccessToApisAllowed bool      `json:"studioAccessToApisAllowed"`
		UniverseAvatarType        string    `json:"universeAvatarType"`
		Genre                     string    `json:"genre"`
	} `json:"data"`
}

// Do makes a request with the provided Request data and returns either a Response or an error
func Do(request Request) (Response, error) {
	var response Response

	v, err := query.Values(request)
	if err != nil {
		return response, err
	}

	req, err := http.NewRequest("GET", "https://games.roblox.com/v1/games?"+v.Encode(), nil)
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
