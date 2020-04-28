package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

//Request ...
type Request struct {
	Query    string
	RegionID int
}

type response struct {
	Suggestions []struct {
		Uri string `json:"Uri"`
	}
}

const searchURL string = "https://avto.pro/"
const requestTimeout time.Duration = time.Second * 120

//TODO: make JSON
const cookie string = `{"IsOpt":false,"Retail":true,"Original":true,"Analog":false,"Used":true,"DeliveryFromCountryISO":"UA","DeliveryFromCityId":null,"DeliveryFromLocationId":690791,"DeliveryToCountryISO":"UA","DeliveryToCityId":null,"DeliveryToLocationId":706483,"LowPrice":true,"ShowCardismantlings":true,"InStock":"all","Currency":"USD","Paging":350,"PartsGrouping":"analogs-and-originals","Sort":"price","Device":0,"UserInShoppingCart":[],"DeliveryToRegionId":null}`

type Page []byte
type Pages [][]byte

func FetchPages(inputUrl string, req Request) (*Pages, error) {
	client := &http.Client{}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal %v: %w", req, err)
	}

	request, err := http.NewRequest("put", inputUrl, bytes.NewReader(data))
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	client.Timeout = requestTimeout
	suggestions, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to make a request to %s: %w", inputUrl, err)
	}
	defer suggestions.Body.Close()

	suggestionsJSON := new(response)
	err = json.NewDecoder(suggestions.Body).Decode(suggestionsJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response %v: %w", suggestions.Body, err)
	}

	pages := make(Pages, len(suggestionsJSON.Suggestions))

	for i, uri := range suggestionsJSON.Suggestions {
		r, err := http.NewRequest("get", searchURL+uri.Uri, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request %s: %w", searchURL+uri.Uri, err)
		}
		r.AddCookie(&http.Cookie{Name: "pref", Value: url.QueryEscape(cookie)})
		response, err := client.Do(r)
		if err != nil {
			return nil, fmt.Errorf("failed to make request %v: %w", r, err)
		}
		page, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body %v: %w", response.Body, err)
		}
		pages[i] = page
	}
	return &pages, err
}
