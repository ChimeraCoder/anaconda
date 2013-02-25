package twitter

import (
    "net/url"
)

func (a TwitterApi) GetHomeTimeline() (timeline []Tweet, err error) {
	v := url.Values{}
	v.Set("include_entities", "true")
	err = a.apiGet("http://api.twitter.com/1.1/statuses/home_timeline.json", v, &timeline)
	return
}

func (a TwitterApi) GetUserTimeline(v url.Values) (timeline []Tweet, err error) {
	err = a.apiGet("http://api.twitter.com/1.1/statuses/user_timeline.json", v, &timeline)
	return
}

func (a TwitterApi) GetMentionsTimeline(v url.Values) (timeline []Tweet, err error) {
	err = a.apiGet("http://api.twitter.com/1.1/statuses/mentions_timeline.json", v, &timeline)
    return
}

func (a TwitterApi) GetRetweetsOfMe(v url.Values) (tweets []Tweet, err error) {
    err = a.apiGet("https://api.twitter.com/1.1/statuses/retweets_of_me.json", v, &tweets)
    return
}


