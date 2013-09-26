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
	a.queryQueue <- query{"https://api.twitter.com/1.1/search/tweets.json", v, &sr, _GET, response_ch}

	timeline = sr.Statuses
	return timeline, (<-response_ch).err
}
