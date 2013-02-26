package twitter

import (
	"net/url"
)

type searchResponse struct {
	Statuses []Tweet
}

func (a TwitterApi) GetSearch(query string, v url.Values) (timeline []Tweet, err error) {
	var sr searchResponse
	v.Set("q", query)
	err = a.apiGet("https://api.twitter.com/1.1/search/tweets.json", v, &sr)
	timeline = sr.Statuses
	return
}
