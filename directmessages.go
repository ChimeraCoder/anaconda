package anaconda

import (
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

func (a TwitterApi) SendDirectMessage(userId string, screenName string, text string) (message DirectMessage, err error) {
	response_ch := make(chan response)

	v := url.Values{}
	v.Set("user_id", userId)
	v.Set("screen_name", screenName)
	v.Set("text", text)

	a.queryQueue <- query{BaseUrl + "/direct_messages/new.json", v, &message, _POST, response_ch}
	return message, (<-response_ch).err
}
