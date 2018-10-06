package anaconda

import (
	"net/url"
	"strconv"
)

func (a TwitterApi) GetUsersLookup(usernames string, v url.Values) (u []User, err error) {
	v = cleanValues(v)
	v.Set("screen_name", usernames)
	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/users/lookup.json", v, &u, _GET, response_ch}
	return u, (<-response_ch).err
}

func (a TwitterApi) GetUsersLookupByIds(ids []int64, v url.Values) (u []User, err error) {
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
	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/users/lookup.json", v, &u, _GET, response_ch}
	return u, (<-response_ch).err
}

func (a TwitterApi) GetUsersShow(username string, v url.Values) (u User, err error) {
	v = cleanValues(v)
	v.Set("screen_name", username)
	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/users/show.json", v, &u, _GET, response_ch}
	return u, (<-response_ch).err
}

func (a TwitterApi) GetUsersShowById(id int64, v url.Values) (u User, err error) {
	v = cleanValues(v)
	v.Set("user_id", strconv.FormatInt(id, 10))
	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/users/show.json", v, &u, _GET, response_ch}
	return u, (<-response_ch).err
}

func (a TwitterApi) GetUserSearch(searchTerm string, v url.Values) (u []User, err error) {
	v = cleanValues(v)
	v.Set("q", searchTerm)
	// Set other values before calling this method:
	// page, count, include_entities
	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/users/search.json", v, &u, _GET, response_ch}
	return u, (<-response_ch).err
}

func (a TwitterApi) GetUsersSuggestions(v url.Values) (c []Category, err error) {
	v = cleanValues(v)
	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/users/suggestions.json", v, &c, _GET, response_ch}
	return c, (<-response_ch).err
}

func (a TwitterApi) GetUsersSuggestionsBySlug(slug string, v url.Values) (s Suggestions, err error) {
	v = cleanValues(v)
	v.Set("slug", slug)
	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/users/suggestions/" + slug + ".json", v, &s, _GET, response_ch}
	return s, (<-response_ch).err
}

// PostUsersReportSpam : Reports and Blocks a User by screen_name
// Reference : https://developer.twitter.com/en/docs/accounts-and-users/mute-block-report-users/api-reference/post-users-report_spam
// If you don't want to block the user you should add
// v.Set("perform_block", "false")
func (a TwitterApi) PostUsersReportSpam(username string, v url.Values) (u User, err error) {
	v = cleanValues(v)
	v.Set("screen_name", username)
	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/users/report_spam.json", v, &u, _POST, response_ch}
	return u, (<-response_ch).err
}

// PostUsersReportSpamById : Reports and Blocks a User by user_id
// Reference : https://developer.twitter.com/en/docs/accounts-and-users/mute-block-report-users/api-reference/post-users-report_spam
// If you don't want to block the user you should add
// v.Set("perform_block", "false")
func (a TwitterApi) PostUsersReportSpamById(id int64, v url.Values) (u User, err error) {
	v = cleanValues(v)
	v.Set("user_id", strconv.FormatInt(id, 10))
	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/users/report_spam.json", v, &u, _POST, response_ch}
	return u, (<-response_ch).err
}

// PostAccountUpdateProfile updates the active users profile with the provided values
func (a TwitterApi) PostAccountUpdateProfile(v url.Values) (u User, err error) {
	v = cleanValues(v)
	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/account/update_profile.json", v, &u, _POST, response_ch}
	return u, (<-response_ch).err
}
