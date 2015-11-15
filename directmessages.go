package anaconda

import (
	"net/url"
	"strconv"
)

func (a TwitterApi) GetDirectMessages(v url.Values) (messages []DirectMessage, err error) {
	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/direct_messages.json", v, &messages, _GET, response_ch}
	return messages, (<-response_ch).err
}

func (a TwitterApi) GetDirectMessagesSent(v url.Values) (messages []DirectMessage, err error) {
	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/direct_messages_sent.json", v, &messages, _GET, response_ch}
	return messages, (<-response_ch).err
}

func (a TwitterApi) GetDirectMessagesShow(v url.Values) (messages []DirectMessage, err error) {
	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/direct_messages/show.json", v, &messages, _GET, response_ch}
	return messages, (<-response_ch).err
}

// https://dev.twitter.com/docs/api/1.1/post/direct_messages/new
func (a TwitterApi) PostDMToScreenName(text, screenName string) (message DirectMessage, err error) {
	v := url.Values{}
	v.Set("screen_name", screenName)
	v.Set("text", text)
	return a.postDirectMessagesImpl(v)
}

// https://dev.twitter.com/docs/api/1.1/post/direct_messages/new
func (a TwitterApi) PostDMToUserId(text string, userId int64) (message DirectMessage, err error) {
	v := url.Values{}
	v.Set("user_id", strconv.FormatInt(userId, 10))
	v.Set("text", text)
	return a.postDirectMessagesImpl(v)
}

func (a TwitterApi) postDirectMessagesImpl(v url.Values) (message DirectMessage, err error) {
	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/direct_messages/new.json", v, &message, _POST, response_ch}
	return message, (<-response_ch).err
}
