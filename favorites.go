package anaconda

import (
	"net/url"
)

func (a TwitterApi) GetFavorites(v url.Values) (favorites []Tweet, err error) {
	response_ch := make(chan response)
	a.queryQueue <- query{BaseUrl + "/favorites/list.json", v, &favorites, _GET, response_ch}
	return favorites, (<-response_ch).err
}

func (a TwitterApi) DeleteFavorite(favorite Tweet) (t Tweet, err error) {
	response_ch := make(chan response)
	a.queryQueue <- query{BaseUrl + "/favorites/destroy.json?id=" + favorite.IdStr, nil, &favorite, _POST, response_ch}
	return favorite, (<-response_ch).err
}
