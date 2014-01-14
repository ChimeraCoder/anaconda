package anaconda

import (
	"net/url"
)

type searchResponse struct {
	Statuses []Tweet
}

func (a TwitterApi) GetSearch(queryString string, v url.Values) (timeline []Tweet, err error) {
	var sr searchResponse

	v = cleanValues(v)
	v.Set("q", queryString)

	response_ch := make(chan response)
	a.queryQueue <- query{BaseUrl + "/search/tweets.json", v, &sr, _GET, response_ch}

	// We have to read from the response channel before assigning to timeline
	// Otherwise this will happen before the responses have been written
	resp := <-response_ch
	err = resp.err
	timeline = sr.Statuses
	return timeline, err
}
