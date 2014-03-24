package anaconda

import (
	"net/url"
	"strconv"
)

func (a TwitterApi) CreateList(name, description string) (list List, err error) {
	v := url.Values{}
	v.Set("name", name)
	v.Set("description", description)

	response_ch := make(chan response)
	a.queryQueue <- query{BaseUrl + "/lists/create.json", v, &list, _POST, response_ch}
	return list, (<-response_ch).err
}

func (a TwitterApi) AddUserToList(screenName string, listID int64) (users []User, err error) {
	v := url.Values{}
	v.Set("list_id", strconv.FormatInt(listID, 10))
	v.Set("screen_name", screenName)

	var addUserToListResponse AddUserToListResponse

	response_ch := make(chan response)
	a.queryQueue <- query{BaseUrl + "/lists/members/create.json", v, &addUserToListResponse, _POST, response_ch}
	return addUserToListResponse.Users, (<-response_ch).err
}

func (a TwitterApi) GetListsOwnedBy(userID int64, count int) (lists []List, err error) {
	v := url.Values{}
	v.Set("user_id", strconv.FormatInt(userID, 10))
	v.Set("count", strconv.FormatInt(userID, 10))

	var listResponse ListResponse

	response_ch := make(chan response)
	a.queryQueue <- query{BaseUrl + "/lists/ownerships.json", v, &listResponse, _GET, response_ch}
	return listResponse.Lists, (<-response_ch).err
}

func (a TwitterApi) GetListTweets(listID int64, includeRTs bool) (tweets []Tweet, err error) {
	v := url.Values{}
	v.Set("list_id", strconv.FormatInt(listID, 10))
	v.Set("include_rts", strconv.FormatBool(includeRTs))

	response_ch := make(chan response)
	a.queryQueue <- query{BaseUrl + "/lists/statuses.json", v, &tweets, _GET, response_ch}
	return tweets, (<-response_ch).err
}
