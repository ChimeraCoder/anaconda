package twitter

import (
	"encoding/json"
	"fmt"
    "strconv"
	"github.com/garyburd/go-oauth/oauth"
	"io/ioutil"
	"net/http"
	"net/url"
)

var OauthClient = oauth.Client{
	TemporaryCredentialRequestURI: "http://api.twitter.com/oauth/request_token",
	ResourceOwnerAuthorizationURI: "http://api.twitter.com/oauth/authenticate",
	TokenRequestURI:               "http://api.twitter.com/oauth/access_token",
}

func SetCredentials(client_marshalled string) (err error) {
	err = json.Unmarshal([]byte(client_marshalled), &OauthClient)
	return
}

// apiGet issues a GET request to the Twitter API and decodes the response JSON to data.
func apiGet(cred *oauth.Credentials, urlStr string, form url.Values, data interface{}) error {
	resp, err := OauthClient.Get(http.DefaultClient, cred, urlStr, form)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return decodeResponse(resp, data)
}

// apiPost issues a POST request to the Twitter API and decodes the response JSON to data.
func apiPost(cred *oauth.Credentials, urlStr string, form url.Values, data interface{}) error {
	resp, err := OauthClient.Post(http.DefaultClient, cred, urlStr, form)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return decodeResponse(resp, data)
}

// decodeResponse decodes the JSON response from the Twitter API.
func decodeResponse(resp *http.Response, data interface{}) error {
	if resp.StatusCode != 200 {
		p, _ := ioutil.ReadAll(resp.Body)

		return fmt.Errorf("Get %s returned status %d, %s", resp.Request.URL, resp.StatusCode, p)
	}
	return json.NewDecoder(resp.Body).Decode(data)
}

func GetHomeTimeline(client *oauth.Credentials) (timeline []Tweet, err error) {
	v := url.Values{}
	v.Set("include_entities", "true")
	err = apiGet(client, "http://api.twitter.com/1.1/statuses/home_timeline.json", v, &timeline)
	return
}

func GetUserTimeline(client *oauth.Credentials, v url.Values) (timeline []Tweet, err error) {
	err = apiGet(client, "http://api.twitter.com/1.1/statuses/user_timeline.json", v, &timeline)
	return
}

func GetMentionsTimeline(client *oauth.Credentials, v url.Values) (timeline []Tweet, err error) {
	err = apiGet(client, "http://api.twitter.com/1.1/statuses/mentions_timeline.json", v, &timeline)
    return
}

func GetRetweetsOfMe(client *oauth.Credentials, v url.Values) (tweets []Tweet, err error) {
    err = apiGet(client, "https://api.twitter.com/1.1/statuses/retweets_of_me.json", v, &tweets)
    return
}

func GetTweet(client *oauth.Credentialis, id int64, v url.Values) (tweet Tweet, err error){
    v.Set("id", strconv.FormatInt(id,10))
    err = apiGet(client, "https://api.twitter.com/1.1/statuses/show.json", v, &tweet)
    return
}

