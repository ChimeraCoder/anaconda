package anaconda

import (
	"fmt"
	"net/url"
)

type Suggestion struct {
	Name string
	Slug string
	Size int64
}

type SuggestionUserList struct {
	Suggestion
	Users []User
}

func (a TwitterApi) GetSuggestions(v url.Values) (s []Suggestion, err error) {
	v = cleanValues(v)
	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/users/suggestions.json", v, &s, _GET, response_ch}
	return s, (<-response_ch).err
}

func (a TwitterApi) GetSuggestionUserList(slug string, v url.Values) (s SuggestionUserList, err error) {
	v = cleanValues(v)
	v.Set("slug", slug)
	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + fmt.Sprintf("/users/suggestions/%s.json", slug), v, &s, _GET, response_ch}
	return s, (<-response_ch).err
}
