package anaconda

import (
	"net/url"
	"strconv"
)

type Cursor struct {
	Previous_cursor     int64
	Previous_cursor_str string

	Ids []int64

	Next_cursor     int64
	Next_cursor_str string
}

type UserCursor struct {
	Previous_cursor     int64
	Previous_cursor_str string
	Next_cursor         int64
	Next_cursor_str     string
	Users               []User
}

type FriendsIdsCursor struct {
	Previous_cursor     int64
	Previous_cursor_str string
	Next_cursor         int64
	Next_cursor_str     string
	Ids                 []int64
}

type FriendsIdsPage struct {
	Ids   []int64
	Error error
}

type Friendship struct {
	Name        string
	Id_str      string
	Id          int64
	Connections []string
	Screen_name string
}

type FollowersPage struct {
	Followers []User
	Error     error
}

//GetFriendshipsNoRetweets s a collection of user_ids that the currently authenticated user does not want to receive retweets from.
//It does not currently support the stringify_ids parameter
func (a TwitterApi) GetFriendshipsNoRetweets() (ids []int64, err error) {
	response_ch := make(chan response)
	a.queryQueue <- query{BaseUrl + "/friendships/no_retweets/ids.json", nil, &ids, _GET, response_ch}
	return ids, (<-response_ch).err
}

func (a TwitterApi) GetFollowersIds(v url.Values) (c Cursor, err error) {
	err = a.apiGet(BaseUrl+"/followers/ids.json", v, &c)
	return
}

func (a TwitterApi) GetFriendsIds(v url.Values) (c Cursor, err error) {
	response_ch := make(chan response)
	a.queryQueue <- query{BaseUrl + "/friends/ids.json", v, &c, _GET, response_ch}
	return c, (<-response_ch).err
}

func (a TwitterApi) GetFriendshipsLookup(v url.Values) (friendships []Friendship, err error) {
	response_ch := make(chan response)
	a.queryQueue <- query{BaseUrl + "/friendships/lookup.json", v, &friendships, _GET, response_ch}
	return friendships, (<-response_ch).err
}

func (a TwitterApi) GetFriendshipsIncoming(v url.Values) (c Cursor, err error) {
	response_ch := make(chan response)
	a.queryQueue <- query{BaseUrl + "/friendships/incoming.json", v, &c, _GET, response_ch}
	return c, (<-response_ch).err
}

func (a TwitterApi) GetFriendshipsOutgoing(v url.Values) (c Cursor, err error) {
	response_ch := make(chan response)
	a.queryQueue <- query{BaseUrl + "/friendships/outgoing.json", v, &c, _GET, response_ch}
	return c, (<-response_ch).err
}

func (a TwitterApi) GetFollowersList(v url.Values) (c UserCursor, err error) {
	response_ch := make(chan response)
	a.queryQueue <- query{BaseUrl + "/followers/list.json", v, &c, _GET, response_ch}
	return c, (<-response_ch).err
}

func (a TwitterApi) GetFriendsList(v url.Values) (c UserCursor, err error) {
	response_ch := make(chan response)
	a.queryQueue <- query{BaseUrl + "/friends/list.json", v, &c, _GET, response_ch}
	return c, (<-response_ch).err
}

// Like GetFollowersList, but returns a channel instead of a cursor and pre-fetches the remaining results
// This channel is closed once all values have been fetched
func (a TwitterApi) GetFollowersListAll(v url.Values) (result chan FollowersPage) {

	result = make(chan FollowersPage)

	if v == nil {
		v = url.Values{}
	}
	go func(a TwitterApi, v url.Values, result chan FollowersPage) {
		// Cursor defaults to the first page ("-1")
		next_cursor := "-1"
		for {
			v.Set("cursor", next_cursor)
			c, err := a.GetFollowersList(v)

			// throttledQuery() handles all rate-limiting errors
			// if GetFollowersList() returns an error, it must be a different kind of error

			result <- FollowersPage{c.Users, err}

			next_cursor = c.Next_cursor_str
			if next_cursor == "0" {
				close(result)
				break
			}
		}
	}(a, v, result)
	return result
}

func (a TwitterApi) GetFollowersUser(id int64, v url.Values) (c Cursor, err error) {
	v = cleanValues(v)
	v.Set("user_id", strconv.FormatInt(id, 10))
	response_ch := make(chan response)
	a.queryQueue <- query{BaseUrl + "/followers/ids.json", v, &c, _GET, response_ch}
	return c, (<-response_ch).err
}

// Like GetFriendsIds, but returns a channel instead of a cursor and pre-fetches the remaining results
// This channel is closed once all values have been fetched
func (a TwitterApi) GetFriendsIdsAll(v url.Values) (result chan FriendsIdsPage) {

	result = make(chan FriendsIdsPage)

	if v == nil {
		v = url.Values{}
	}
	go func(a TwitterApi, v url.Values, result chan FriendsIdsPage) {
		// Cursor defaults to the first page ("-1")
		next_cursor := "-1"
		for {
			v.Set("cursor", next_cursor)
			c, err := a.GetFriendsIds(v)

			// throttledQuery() handles all rate-limiting errors
			// if GetFollowersList() returns an error, it must be a different kind of error

			result <- FriendsIdsPage{c.Ids, err}

			next_cursor = c.Next_cursor_str
			if next_cursor == "0" {
				close(result)
				break
			}
		}
	}(a, v, result)
	return result
}

func (a TwitterApi) GetFriendsUser(id int64, v url.Values) (c Cursor, err error) {
	v = cleanValues(v)
	v.Set("user_id", strconv.FormatInt(id, 10))
	response_ch := make(chan response)
	a.queryQueue <- query{BaseUrl + "/friends/ids.json", v, &c, _GET, response_ch}
	return c, (<-response_ch).err
}

// FollowUserId follows the user with the specified userId.
// This implements the /friendships/create endpoint, though the function name
// uses the terminology 'follow' as this is most consistent with colloquial Twitter terminology.
func (a TwitterApi) FollowUserId(userId int64, v url.Values) (user User, err error) {
	v = cleanValues(v)
	v.Set("user_id", strconv.FormatInt(userId, 10))
	return a.postFriendshipsCreateImpl(v)
}

// FollowUserId follows the user with the specified screenname (username).
// This implements the /friendships/create endpoint, though the function name
// uses the terminology 'follow' as this is most consistent with colloquial Twitter terminology.
func (a TwitterApi) FollowUser(screenName string) (user User, err error) {
	v := url.Values{}
	v.Set("screen_name", screenName)
	return a.postFriendshipsCreateImpl(v)
}

func (a TwitterApi) postFriendshipsCreateImpl(v url.Values) (user User, err error) {
	response_ch := make(chan response)
	a.queryQueue <- query{BaseUrl + "/friendships/create.json", v, &user, _POST, response_ch}
	return user, (<-response_ch).err
}

// UnfollowUserId unfollows the user with the specified userId.
// This implements the /friendships/destroy endpoint, though the function name
// uses the terminology 'unfollow' as this is most consistent with colloquial Twitter terminology.
func (a TwitterApi) UnfollowUserId(userId int64) (u User, err error) {
	v := url.Values{}
	v.Set("user_id", strconv.FormatInt(userId, 10))
	// Set other values before calling this method:
	// page, count, include_entities
	response_ch := make(chan response)
	a.queryQueue <- query{BaseUrl + "/friendships/destroy.json", v, &u, _POST, response_ch}
	return u, (<-response_ch).err
}

// UnfollowUser unfollows the user with the specified screenname (username)
// This implements the /friendships/destroy endpoint, though the function name
// uses the terminology 'unfollow' as this is most consistent with colloquial Twitter terminology.
func (a TwitterApi) UnfollowUser(screenname string) (u User, err error) {
	v := url.Values{}
	v.Set("screen_name", screenname)
	// Set other values before calling this method:
	// page, count, include_entities
	response_ch := make(chan response)
	a.queryQueue <- query{BaseUrl + "/friendships/destroy.json", v, &u, _POST, response_ch}
	return u, (<-response_ch).err
}
