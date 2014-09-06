package anaconda

import (
	"net/url"
)

func (a TwitterApi) GetFavorites(v url.Values) (favorites []Tweet, err error) {
	response_ch := make(chan response)
	a.queryQueue <- query{BaseUrl + "/favorites/list.json", v, &favorites, _GET, response_ch}
	return favorites, (<-response_ch).err
}

func (a TwitterApi) DeleteFavorite(id int64) error {
	v := url.Values{}
	v.Set("id", string(id))
	response_ch := make(chan response)
	a.queryQueue <- query{BaseUrl + "favorites/destroy.json", v, nil, _POST, response_ch}
	return (<-response_ch).err
}
