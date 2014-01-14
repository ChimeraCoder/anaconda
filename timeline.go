package anaconda

import (
	"net/url"
)

func (a TwitterApi) GetHomeTimeline() (timeline []Tweet, err error) {
	v := url.Values{}
	v.Set("include_entities", "true")

	response_ch := make(chan response)
	a.queryQueue <- query{BaseUrl + "/statuses/home_timeline.json", v, &timeline, _GET, response_ch}
	return timeline, (<-response_ch).err
}

func (a TwitterApi) GetUserTimeline(v url.Values) (timeline []Tweet, err error) {
	response_ch := make(chan response)
	a.queryQueue <- query{BaseUrl + "/statuses/user_timeline.json", v, &timeline, _GET, response_ch}
	return timeline, (<-response_ch).err
}

func (a TwitterApi) GetMentionsTimeline(v url.Values) (timeline []Tweet, err error) {
	response_ch := make(chan response)
	a.queryQueue <- query{BaseUrl + "/statuses/mentions_timeline.json", v, &timeline, _GET, response_ch}
	return timeline, (<-response_ch).err
}

func (a TwitterApi) GetRetweetsOfMe(v url.Values) (tweets []Tweet, err error) {
	response_ch := make(chan response)
	a.queryQueue <- query{BaseUrl + "/statuses/retweets_of_me.json", v, &tweets, _GET, response_ch}
	return tweets, (<-response_ch).err
}
