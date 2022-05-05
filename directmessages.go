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
	a.queryQueue <- query{a.baseUrl + "/direct_messages/sent.json", v, &messages, _GET, response_ch}
	return messages, (<-response_ch).err
}

func (a TwitterApi) GetDirectMessagesShow(v url.Values) (message DirectMessage, err error) {
	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/direct_messages/show.json", v, &message, _GET, response_ch}
	return message, (<-response_ch).err
}

// https://developer.twitter.com/en/docs/direct-messages/sending-and-receiving/api-reference/new-message
func (a TwitterApi) PostDMToScreenName(text, screenName string) (message DirectMessage, err error) {
	v := url.Values{}
	v.Set("screen_name", screenName)
	v.Set("text", text)
	return a.postDirectMessagesImpl(v)
}

// https://developer.twitter.com/en/docs/direct-messages/sending-and-receiving/api-reference/new-message
func (a TwitterApi) PostDMToUserId(text string, userId int64) (message DirectMessage, err error) {
	v := url.Values{}
	v.Set("user_id", strconv.FormatInt(userId, 10))
	v.Set("text", text)
	return a.postDirectMessagesImpl(v)
}

// DeleteDirectMessage will destroy (delete) the direct message with the specified ID.
// https://developer.twitter.com/en/docs/direct-messages/sending-and-receiving/api-reference/delete-message
func (a TwitterApi) DeleteDirectMessage(id int64, includeEntities bool) (message DirectMessage, err error) {
	v := url.Values{}
	v.Set("id", strconv.FormatInt(id, 10))
	v.Set("include_entities", strconv.FormatBool(includeEntities))
	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/direct_messages/destroy.json", v, &message, _POST, response_ch}
	return message, (<-response_ch).err
}

func (a TwitterApi) postDirectMessagesImpl(v url.Values) (message DirectMessage, err error) {
	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/direct_messages/new.json", v, &message, _POST, response_ch}
	return message, (<-response_ch).err
}

// IndicateTyping will create a typing indicator
// https://developer.twitter.com/en/docs/direct-messages/typing-indicator-and-read-receipts/api-reference/new-typing-indicator
func (a TwitterApi) IndicateTyping(id int64) (err error) {
	v := url.Values{}
	v.Set("recipient_id", strconv.FormatInt(id, 10))
	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/direct_messages/indicate_typing.json", v, nil, _POST, response_ch}
	return (<-response_ch).err
}
