package anaconda

import (
	"net/url"
)

func (a TwitterApi) GetUsersLookup(usernames string, v url.Values) (u []TwitterUser, err error) {
	v = cleanValues(v)
	v.Set("screen_name", usernames)
	//err = a.apiGet("http://api.twitter.com/1.1/users/lookup.json", v, &u)
	response_ch := make(chan response)
	a.queryQueue <- query{"http://api.twitter.com/1.1/users/lookup.json", v, &u, _GET, response_ch}
	response := <-response_ch
	return u, response.err
}
