package anaconda

import (
	"net/url"
)

func (a TwitterApi) GetUsersLookup(usernames string, v url.Values) (u []User, err error) {
	v = cleanValues(v)
	v.Set("screen_name", usernames)
	response_ch := make(chan response)
	a.queryQueue <- query{"http://api.twitter.com/1.1/users/lookup.json", v, &u, _GET, response_ch}
	return u, (<-response_ch).err
}
