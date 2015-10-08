package anaconda

import "net/url"

type Location struct {
	Name  string `json:"name"`
	Woeid int    `json:"woeid"`
}

type Trend struct {
	Name            string `json:"name"`
	Query           string `json:"query"`
	Url             string `json:"url"`
	PromotedContent string `json:"promoted_content"`
}

type TrendResponse struct {
	Trends    []Trend    `json:"trends"`
	AsOf      string     `json:"as_of"`
	CreatedAt string     `json:"created_at"`
	Locations []Location `json:"locations"`
}

type TrendLocation struct {
	Country     string `json:"country"`
	CountryCode string `json:"countryCode"`
	Name        string `json:"name"`
	ParentId    int    `json:"parentid"`
	PlaceType   struct {
		Code int    `json:"code"`
		Name string `json:"name"`
	} `json:"placeType"`
	Url   string `json:"url"`
	Woeid int32  `json:"woeid"`
}

// https://dev.twitter.com/rest/reference/get/trends/place
func (a TwitterApi) GetTrendsByPlace(v url.Values) (trendResp TrendResponse, err error) {
	trendResponse := TrendResponse{}
	response_ch := make(chan response)
	a.queryQueue <- query{BaseUrl + "/trends/place.json", v, &[]interface{}{&trendResponse}, _GET, response_ch}
	return trendResponse, (<-response_ch).err
}

// https://dev.twitter.com/rest/reference/get/trends/available
func (a TwitterApi) GetTrendsAvailableLocations(v url.Values) (locations []TrendLocation, err error) {
	response_ch := make(chan response)
	a.queryQueue <- query{BaseUrl + "/trends/available.json", v, &locations, _GET, response_ch}
	return locations, (<-response_ch).err
}

// https://dev.twitter.com/rest/reference/get/trends/closest
func (a TwitterApi) GetTrendsClosestLocations(v url.Values) (locations []TrendLocation, err error) {
	response_ch := make(chan response)
	a.queryQueue <- query{BaseUrl + "/trends/closest.json", v, &locations, _GET, response_ch}
	return locations, (<-response_ch).err
}
