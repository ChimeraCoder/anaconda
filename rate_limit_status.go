package anaconda

import "net/url"
import "fmt"

type RateLimitSearchResult struct {
	RateLimitContext struct {
		AccessToken string `json:"access_token"`
	} `json:"rate_limit_context"`
	Resources struct {
		Geo struct {
			Similar_places struct {
				Limit     int `json:"limit"`
				Remaining int `json:"remaining"`
				Reset     int `json:"reset"`
			} `json:"/geo/similar_places"`
			Place_id struct {
				Limit     int `json:"limit"`
				Remaining int `json:"remaining"`
				Reset     int `json:"reset"`
			} `json:"/geo/id/:place_id"`
			Reverse_geocode struct {
				Limit     int `json:"limit"`
				Remaining int `json:"remaining"`
				Reset     int `json:"reset"`
			} `json:"/geo/reverse_geocode"`
			Search struct {
				Limit     int `json:"limit"`
				Remaining int `json:"remaining"`
				Reset     int `json:"reset"`
			} `json:"/geo/search"`
		} `json:"geo"`
	} `json:"resources"`
}

func (a TwitterApi) GetRateLimitStatus(v url.Values) (r RateLimitSearchResult, err error) {
	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/application/rate_limit_status.json", v, &r, _GET, response_ch}
	return r, (<-response_ch).err
}
