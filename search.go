package anaconda

import (
	"net/url"
)

type SearchMetadata struct {
	CompletedIn   float32 `json:"completed_in"`
	MaxId         int64   `json:"max_id"`
	MaxIdString   string  `json:"max_id_str"`
	Query         string  `json:"query"`
	RefreshUrl    string  `json:"refresh_url"`
	Count         int     `json:"count"`
	SinceId       int64   `json:"since_id"`
	SinceIdString string  `json:"since_id_str"`
}

type SearchResponse struct {
	Statuses []Tweet        `json:"statuses"`
	Metadata SearchMetadata `json:"search_metadata"`
}

func (a TwitterApi) GetSearch(queryString string, v url.Values) (sr SearchResponse, err error) {
	v = cleanValues(v)
	v.Set("q", queryString)

	response_ch := make(chan response)
	a.queryQueue <- query{BaseUrl + "/search/tweets.json", v, &sr, _GET, response_ch}

	// We have to read from the response channel before assigning to timeline
	// Otherwise this will happen before the responses have been written
	resp := <-response_ch
	err = resp.err
	return sr, err
}
