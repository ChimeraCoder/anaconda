package anaconda

import (
	"net/http"
	"net/url"
)

type OEmbed struct {
	Type          string
	Width         int
	Cache_age     string
	Height        int
	Author_url    string
	Html          string
	Version       string
	Provider_name string
	Provider_url  string
	Url           string
	Author_name   string
}

// No authorization on this endpoint. Its the only one.
func (a TwitterApi) GetOEmbed(v url.Values) (o OEmbed, err error) {
	resp, err := http.Get("http://api.twitter.com/1/statuses/oembed.json?" + v.Encode())
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = decodeResponse(resp, &o)
	return
}
