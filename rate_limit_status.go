package anaconda

import (
	"net/url"
	"strings"
)

type RateLimitStatusResponse struct {
	RateLimitContext RateLimitContext                   `json:"rate_limit_context"`
	Resources        map[string]map[string]BaseResource `json:"resources"`
}

type RateLimitContext struct {
	AccessToken string `json:"access_token"`
}

type BaseResource struct {
	Limit     int `json:"limit"`
	Remaining int `json:"remaining"`
	Reset     int `json:"reset"`
}

func (a TwitterApi) GetRateLimits(r []string) (rateLimitStatusResponse RateLimitStatusResponse, err error) {
	resources := strings.Join(r, ",")
	v := url.Values{}
	v.Set("resources", resources)
	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/application/rate_limit_status.json", v, &rateLimitStatusResponse, _GET, response_ch}
	return rateLimitStatusResponse, (<-response_ch).err
}
