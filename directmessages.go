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
