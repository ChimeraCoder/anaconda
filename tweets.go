package anaconda

import (
	"fmt"
	"net/url"
	"strconv"
)

func (a TwitterApi) GetTweet(id int64, v url.Values) (tweet Tweet, err error) {
	v.Set("id", strconv.FormatInt(id, 10))
	err = a.apiGet("https://api.twitter.com/1.1/statuses/show.json", v, &tweet)
	return
}

func (a TwitterApi) GetRetweets(id int64, v url.Values) (tweets []Tweet, err error) {
	err = a.apiGet(fmt.Sprintf("https://api.twitter.com/1.1/statuses/retweets/%d.json", id), v, &tweets)
	return
}

//PostTweet will create a tweet with the specified status message
func (a TwitterApi) PostTweet(status string, v url.Values) (tweet Tweet, err error) {
	v.Set("status", status)
	err = a.apiPost("https://api.twitter.com/1.1/statuses/update.json", v, &tweet)
	return
}

//DeleteTweet will destroy (delete) the status (tweet) with the specified ID, assuming that the authenticated user is the author of the status (tweet).
//If trimUser is set to true, only the user's Id will be provided in the user object returned.
func (a TwitterApi) DeleteTweet(id int64, trimUser bool) (tweet Tweet, err error) {
	v := url.Values{}
	if trimUser {
		v.Set("trim_user", "t")
	}
	err = a.apiPost(fmt.Sprintf("https://api.twitter.com/1.1/statuses/destroy/%d.json", id), v, &tweet)
	return
}

//Retweet will retweet the status (tweet) with the specified ID.
//trimUser functions as in DeleteTweet
func (a TwitterApi) Retweet(id int64, trimUser bool) (rt Retweet, err error) {
	v := url.Values{}
	if trimUser {
		v.Set("trim_user", "t")
	}
	err = a.apiPost(fmt.Sprintf("https://api.twitter.com/1.1/statuses/retweet/%d.json", id), v, &rt)
	return
}
