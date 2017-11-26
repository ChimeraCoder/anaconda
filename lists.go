package anaconda

import (
	"net/url"
	"strconv"
	"strings"
)

// CreateList implements /lists/create.json
func (a TwitterApi) CreateList(name, description string, v url.Values) (list List, err error) {
	v = cleanValues(v)
	v.Set("name", name)
	v.Set("description", description)

	ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/lists/create.json", v, &list, _POST, ch}
	return list, (<-ch).err
}

// AddUserToList implements /lists/members/create.json
func (a TwitterApi) AddUserToList(screenName string, listID int64, v url.Values) (users []User, err error) {
	v = cleanValues(v)
	v.Set("list_id", strconv.FormatInt(listID, 10))
	v.Set("screen_name", screenName)

	var addUserToListResponse AddUserToListResponse

	ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/lists/members/create.json", v, &addUserToListResponse, _POST, ch}
	return addUserToListResponse.Users, (<-ch).err
}

// AddMultipleUsersToList implements /lists/members/create_all.json
func (a TwitterApi) AddMultipleUsersToList(screenNames []string, listID int64, v url.Values) (list List, err error) {
	v = cleanValues(v)
	v.Set("list_id", strconv.FormatInt(listID, 10))
	v.Set("screen_name", strings.Join(screenNames, ","))

	ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/lists/members/create_all.json", v, &list, _POST, ch}
	r := <-ch
	return list, r.err
}

// GetListsOwnedBy implements /lists/ownerships.json
// screen_name, count, and cursor are all optional values
func (a TwitterApi) GetListsOwnedBy(userID int64, v url.Values) (lists []List, err error) {
	v = cleanValues(v)
	v.Set("user_id", strconv.FormatInt(userID, 10))

	var listResponse ListResponse

	ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/lists/ownerships.json", v, &listResponse, _GET, ch}
	return listResponse.Lists, (<-ch).err
}

// GetListTweets implements /lists/statuses.json
// Returns all tweets from users in a specific list
func (a TwitterApi) GetListTweets(listID int64, includeRTs bool, v url.Values) (tweets []Tweet, err error) {
	v = cleanValues(v)
	v.Set("list_id", strconv.FormatInt(listID, 10))
	v.Set("include_rts", strconv.FormatBool(includeRTs))

	ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/lists/statuses.json", v, &tweets, _GET, ch}
	return tweets, (<-ch).err
}

// GetList implements /lists/show.json
func (a TwitterApi) GetList(listID int64, v url.Values) (list List, err error) {
	v = cleanValues(v)
	v.Set("list_id", strconv.FormatInt(listID, 10))

	ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/lists/show.json", v, &list, _GET, ch}
	return list, (<-ch).err
}

func (a TwitterApi) GetListTweetsBySlug(slug string, ownerScreenName string, includeRTs bool, v url.Values) (tweets []Tweet, err error) {
	v = cleanValues(v)
	v.Set("slug", slug)
	v.Set("owner_screen_name", ownerScreenName)
	v.Set("include_rts", strconv.FormatBool(includeRTs))

	ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/lists/statuses.json", v, &tweets, _GET, ch}
	return tweets, (<-ch).err
}
