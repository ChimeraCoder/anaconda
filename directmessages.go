package twitter

import (
	"net/url"
)

func (a TwitterApi) GetDirectMessages(v url.Values) (messages []DirectMessage, err error) {
	err = a.apiGet("https://api.twitter.com/1.1/direct_messages.json", v, &messages)
	return
}

func (a TwitterApi) GetDirectMessagesSent(v url.Values) (messages []DirectMessage, err error) {
	err = a.apiGet("https://api.twitter.com/1.1/direct_messages.json", v, &messages)
	return
}

func (a TwitterApi) GetDirectMessagesShow(v url.Values) (messages []DirectMessage, err error) {
	err = a.apiGet("https://api.twitter.com/1.1/direct_messages/sent.json", v, &messages)
	return
}
