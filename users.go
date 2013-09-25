package anaconda

import (
	"net/url"
	"strconv"
)

func (a TwitterApi) GetUsersLookup(usernames string, v url.Values) (u []TwitterUser, err error) {
	v = cleanValues(v)
	v.Set("screen_name", usernames)
	response_ch := make(chan response)
	a.queryQueue <- query{"http://api.twitter.com/1.1/users/lookup.json", v, &u, _GET, response_ch}
	return u, (<-response_ch).err
}

func (a TwitterApi) GetUsersLookupByIds(ids []int64, v url.Values) (u []TwitterUser, err error) {
	var pids string
	for w, i := range ids {
		//pids += strconv.Itoa(i)
		pids += strconv.FormatInt(i, 10)
		if w != len(ids)-1 {
			pids += ","
		}
	}
	v = cleanValues(v)
	v.Set("user_id", pids)
	err = a.apiGet("http://api.twitter.com/1.1/users/lookup.json", v, &u)
	return
}
