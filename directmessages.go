package anaconda

import (
	"fmt"
	"net/url"
)

func (a TwitterApi) GetDirectMessages(v url.Values) (messages []DirectMessage, err error) {
	response_ch := make(chan response)
	a.queryQueue <- query{BaseUrl + "/direct_messages.json", v, &messages, _GET, response_ch}
	return messages, (<-response_ch).err
}

func (a TwitterApi) GetDirectMessagesSent(v url.Values) (messages []DirectMessage, err error) {
	response_ch := make(chan response)
	a.queryQueue <- query{BaseUrl + "/direct_messages_sent.json", v, &messages, _GET, response_ch}
	return messages, (<-response_ch).err
}

func (a TwitterApi) GetDirectMessagesShow(v url.Values) (messages []DirectMessage, err error) {
	response_ch := make(chan response)
	a.queryQueue <- query{BaseUrl + "/direct_messages/show.json", v, &messages, _GET, response_ch}
	return messages, (<-response_ch).err
}

func (a TwitterApi) SendDirectMessageUserId(userId int64, text string) (message DirectMessage, err error) {
	v := url.Values{}
	v.Set("user_id", fmt.Sprint(userId))

	return a.sendDirectMessage(text, v)
}

func (a TwitterApi) SendDirectMessageScreenName(screenName string, text string) (message DirectMessage, err error) {
	v := url.Values{}
	v.Set("screen_name", screenName)

	return a.sendDirectMessage(text, v)
}

func (a TwitterApi) sendDirectMessage(text string, v url.Values) (message DirectMessage, err error) {
	response_ch := make(chan response)
	v.Set("text", text)

	a.queryQueue <- query{BaseUrl + "/direct_messages/new.json", v, &message, _POST, response_ch}
	return message, (<-response_ch).err
}
