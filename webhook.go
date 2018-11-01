package anaconda

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
)

//GetActivityWebhooks represents the twitter account_activity webhook
//Returns all URLs and their statuses for the given app. Currently,
//only one webhook URL can be registered to an application except in the Enterprise API
//https://developer.twitter.com/en/docs/accounts-and-users/subscribe-account-activity/api-reference/aaa-premium#get-account-activity-all-webhooks
func (a TwitterApi) GetActivityWebhooks(v url.Values) (u []WebHookResp, err error) {
	v = cleanValues(v)
	responseCh := make(chan response)
	a.queryQueue <- query{a.activityUrl + "webhooks.json", v, &u, _GET, responseCh}
	return u, (<-responseCh).err
}

//WebHookResp represents the Get webhook responses
type WebHookResp struct {
	ID        string
	URL       string
	Valid     bool
	CreatedAt string
}

type BearerToken struct {
	Type  string `json:"token_type"`
	Token string `json:"access_token"`
}

//SetActivityWebhooks represents to set twitter account_activity webhook
//Registers a new webhook URL for the given application context.
//The URL will be validated via CRC request before saving. In case the validation fails,
//a comprehensive error message is returned to the requester.
//Only one webhook URL can be registered to an application except in the Enterprise API.
//https://developer.twitter.com/en/docs/accounts-and-users/subscribe-account-activity/api-reference/aaa-premium#post-account-activity-all-env-name-webhooks
func (a TwitterApi) SetActivityWebhooks(v url.Values) (u WebHookResp, err error) {
	v = cleanValues(v)
	responseCh := make(chan response)
	a.queryQueue <- query{a.activityUrl + "webhooks.json", v, &u, _POST, responseCh}
	return u, (<-responseCh).err
}

//DeleteActivityWebhooks Removes a webhook from the provided application’s configuration.
//https://developer.twitter.com/en/docs/accounts-and-users/subscribe-account-activity/api-reference/aaa-premium#delete-account-activity-all-env-name-webhooks-webhook-id
func (a TwitterApi) DeleteActivityWebhooks(v url.Values, webhookID string) (u interface{}, err error) {
	v = cleanValues(v)
	responseCh := make(chan response)
	a.queryQueue <- query{a.activityUrl + "webhooks/" + webhookID + ".json", v, &u, _DELETE, responseCh}
	return u, (<-responseCh).err
}

//PutActivityWebhooks Updates a webhook by triggering a challenge response check (CRC) which, if
//successful, sets its status to valid.
//https://developer.twitter.com/en/docs/accounts-and-users/subscribe-account-activity/api-reference/aaa-premium#put-account-activity-all-env-name-webhooks-webhook-id
func (a TwitterApi) PutActivityWebhooks(v url.Values, webhookID string) (u interface{}, err error) {
	v = cleanValues(v)
	responseCh := make(chan response)
	a.queryQueue <- query{a.activityUrl + "webhooks/" + webhookID + ".json", v, &u, _PUT, responseCh}
	return u, (<-responseCh).err
}

//SetWHSubscription Subscribes the provided app to events for the provided user context.
//When subscribed, all events for the provided user will be sent to the app’s webhook via POST request.
//https://developer.twitter.com/en/docs/accounts-and-users/subscribe-account-activity/api-reference/aaa-premium#post-account-activity-all-env-name-subscriptions
func (a TwitterApi) SetWHSubscription(v url.Values, webhookID string) (u interface{}, err error) {
	v = cleanValues(v)
	responseCh := make(chan response)
	if a.env != "" {
		a.queryQueue <- query{a.activityUrl + "subscriptions.json", v, &u, _POST, responseCh}
	} else {
		a.queryQueue <- query{a.activityUrl + "webhooks/" + webhookID + "/subscriptions/all.json", v, &u, _POST, responseCh}
	}
	return u, (<-responseCh).err
}

//GetWHSubscription Determines if a webhook configuration is subscribed to the provided user’s account
//https://developer.twitter.com/en/docs/accounts-and-users/subscribe-account-activity/api-reference/aaa-premium#get-account-activity-all-env-name-subscriptions
func (a TwitterApi) GetWHSubscription(v url.Values, webhookID string) (u interface{}, err error) {
	v = cleanValues(v)
	responseCh := make(chan response)
	if a.env != "" {
		a.queryQueue <- query{a.activityUrl + "subscriptions.json", v, &u, _GET, responseCh}
	} else {
		a.queryQueue <- query{a.activityUrl + "webhooks/" + webhookID + "/subscriptions/all.json", v, &u, _GET, responseCh}
	}
	return u, (<-responseCh).err
}

//ListWHSubscriptions Returns a list of active subscriptions
//https://developer.twitter.com/en/docs/accounts-and-users/subscribe-account-activity/api-reference/aaa-premium#get-account-activity-all-env-name-subscriptions-list
func (a TwitterApi) ListWHSubscriptions(v url.Values) (u interface{}, err error) {
	v = cleanValues(v)
	responseCh := make(chan response)

	a.queryQueue <- query{a.activityUrl + "subscriptions/list.json", v, &u, _GETBEARER, responseCh}
	return u, (<-responseCh).err
}

//DeleteWHSubscription Deactivates subscription for the provided user context and app. After deactivation,
//all events for the requesting user will no longer be sent to the webhook URL.
//https://dev.twitter.com/webhooks/reference/del/account_activity/webhooks
func (a TwitterApi) DeleteWHSubscription(v url.Values, webhookID string) (u interface{}, err error) {
	v = cleanValues(v)
	responseCh := make(chan response)
	if a.env != "" {
		a.queryQueue <- query{a.activityUrl + "subscriptions.json", v, &u, _DELETE, responseCh}
	} else {
		a.queryQueue <- query{a.activityUrl + "webhooks/" + webhookID + "/subscriptions/all.json", v, &u, _DELETE, responseCh}
	}
	return u, (<-responseCh).err
}

//CountWHSubscriptions returns the count of subscriptions active on a given webhook
//https://developer.twitter.com/en/docs/accounts-and-users/subscribe-account-activity/api-reference/aaa-premium#get-account-activity-all-subscriptions-count
func (a TwitterApi) CountWHSubscriptions(v url.Values) (u interface{}, err error) {
	v = cleanValues(v)
	responseCh := make(chan response)
	//note lack of environment name here, even for Premium

	a.queryQueue <- query{a.baseUrl + "/account_activity/all/subscriptions/count.json", v, &u, _GETBEARER, responseCh}
	return u, (<-responseCh).err
}

//RespondCRC responds to a CRC request from Twitter.
//Should be called in response to a GET request to the callback URL.
//https://developer.twitter.com/en/docs/accounts-and-users/subscribe-account-activity/guides/securing-webhooks
func (a TwitterApi) RespondCRC(tok string, w http.ResponseWriter) {
	mac := hmac.New(sha256.New, []byte(a.oauthClient.Credentials.Secret))
	mac.Write([]byte(tok))
	resp := "{ \"response_token\": \"sha256=" + base64.StdEncoding.EncodeToString(mac.Sum(nil)) + "\" }"
	fmt.Fprint(w, resp)
}
