package anaconda

import (
	"fmt"
	"net/url"
	"strconv"
)

func (a TwitterApi) GetUsersLookup(usernames string, v url.Values) (u []TwitterUser, err error) {
	v = cleanValues(v)
	v.Set("screen_name", usernames)
	err = a.apiGet("http://api.twitter.com/1.1/users/lookup.json", v, &u)
	return
}

func (a TwitterApi) GetUsersLookupByIds(ids []int, v url.Values) (u []TwitterUser, err error) {
	var pids string
	for w, i := range ids {
		pids += strconv.Itoa(i)
		if w != len(ids)-1 {
			pids += ","
		}
	}
	v = cleanValues(v)
	v.Set("user_id", pids)
	fmt.Println("Foo!", pids)
	err = a.apiGet("http://api.twitter.com/1.1/users/lookup.json", v, &u)
	return
}
