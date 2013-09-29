package anaconda

import (
	"net/url"
)

type Cursor struct {
	Previous_cursor     int64
	Previous_cursor_str string

	Ids []int64

	Next_cursor     int64
	Next_cursor_str string
}

type TwitterUserCursor struct {
    Previous_cursor int64
    Previous_cursor_str string
    Next_cursor int64
    Next_cursor_str string
    Users []TwitterUser
}

type Friendship struct {
	Name        string
	Id_str      string
	Id          int64
	Connections []string
	Screen_name string
}

//GetFriendshipsNoRetweets returns a collection of user_ids that the currently authenticated user does not want to receive retweets from.
//It does not currently support the stringify_ids parameter
func (a TwitterApi) GetFriendshipsNoRetweets() (ids []int64, err error) {
	err = a.apiGet("https://api.twitter.com/1.1/friendships/no_retweets/ids.json", nil, &ids)
	return
}

func (a TwitterApi) GetFollowersIds(v url.Values) (c Cursor, err error) {
	err = a.apiGet("https://api.twitter.com/1.1/followers/ids.json", v, &c)
	return
}

func (a TwitterApi) GetFriendsIds(v url.Values) (c Cursor, err error) {
	err = a.apiGet("https://api.twitter.com/1.1/friends/ids.json", v, &c)
	return
}

func (a TwitterApi) GetFriendshipsLookup(v url.Values) (friendships []Friendship, err error) {
	err = a.apiGet("http://api.twitter.com/1.1/friendships/lookup.json", v, &friendships)
	return
}

func (a TwitterApi) GetFriendshipsIncoming(v url.Values) (c Cursor, err error) {
	err = a.apiGet("https://api.twitter.com/1.1/friendships/incoming.json", v, &c)
	return
}

func (a TwitterApi) GetFriendshipsOutgoing(v url.Values) (c Cursor, err error) {
	err = a.apiGet("http://api.twitter.com/1.1/friendships/outgoing.json", v, &c)
	return
}

func (a TwitterApi) GetFollowersList(v url.Values) (c TwitterUserCursor, err error) {
    err = a.apiGet("https://api.twitter.com/1.1/followers/list.json", v, &c)
	return
}
