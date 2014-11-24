package anaconda

import (
	"net/url"
)

type  statusesResponse struct {
	Statuses []Tweet
}

func (a TwitterApi) GetStatusesMentionsTimeline(v url.Values) (timeline []Tweet, err error) {
	v = cleanValues(v)

	response_ch := make(chan response)
	a.queryQueue <- query{BaseUrl + "/statuses/mentions_timeline.json", v, &timeline, _GET, response_ch}
	return timeline, (<-response_ch).err
}

func (a TwitterApi) GetStatusesUserTimeline(v url.Values) (timeline []Tweet, err error) {
	v = cleanValues(v)

	response_ch := make(chan response)
	a.queryQueue <- query{BaseUrl + "/statuses/user_timeline.json", v, &timeline, _GET, response_ch}
	return timeline, (<-response_ch).err
}
