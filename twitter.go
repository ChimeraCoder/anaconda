//Package anaconda provides structs and functions for accessing version 1.1
//of the Twitter API.
//
//Successful API queries return native Go structs that can be used immediately,
//with no need for type assertions.
//
//Authentication
//
//If you already have the access token (and secret) for your user (Twitter provides this for your own account on the developer portal), creating the client is simple:
//
//  anaconda.SetConsumerKey("your-consumer-key")
//  anaconda.SetConsumerSecret("your-consumer-secret")
//  api := anaconda.NewTwitterApi("your-access-token", "your-access-token-secret")
//
//
//Queries
//
//Executing queries on an authenticated TwitterApi struct is simple.
//
//  searchResult, _ := api.GetSearch("golang", nil)
//  for _ , tweet := range searchResult {
//      fmt.Print(tweet.Text)
//  }
//
//Certain endpoints allow separate optional parameter; if desired, these can be passed as the final parameter. 
//
//  v := url.Values{}
//  v.Set("count", "30")
//  result, err := api.GetSearch("golang", v)
//
//
//Endpoints
//
//Anaconda implements most of the endpoints defined in the Twitter API documentation: https://dev.twitter.com/docs/api/1.1.
//For clarity, in most cases, the function name is simply the name of the HTTP method and the endpoint (e.g., the endpoint `GET /friendships/incoming` is provided by the function `GetFriendshipsIncoming`).
//
//In a few cases, a shortened form has been chosen to make life easier (for example, retweeting is simply the function `Retweet`)
//
//More detailed information about the behavior of each particular endpoint can be found at the official Twitter API documentation.
package anaconda

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/go-oauth/oauth"
	"io/ioutil"
	"net/http"
	"net/url"
)

var oauthClient = oauth.Client{
	TemporaryCredentialRequestURI: "http://api.twitter.com/oauth/request_token",
	ResourceOwnerAuthorizationURI: "http://api.twitter.com/oauth/authenticate",
	TokenRequestURI:               "http://api.twitter.com/oauth/access_token",
}

type TwitterApi struct {
	Credentials *oauth.Credentials
}

//NewTwitterApi takes an user-specific access token and secret and returns a TwitterApi struct for that user.
//The TwitterApi struct can be used for accessing any of the endpoints available.
func NewTwitterApi(access_token string, access_token_secret string) TwitterApi {
	return TwitterApi{&oauth.Credentials{Token: access_token, Secret: access_token_secret}}
}

//SetConsumerKey will set the application-specific consumer_key used in the initial OAuth process
//This key is listed on https://dev.twitter.com/apps/YOUR_APP_ID/show
func SetConsumerKey(consumer_key string) {
	oauthClient.Credentials.Token = consumer_key
}

//SetConsumerSecret will set the application-specific secret used in the initial OAuth process
//This secret is listed on https://dev.twitter.com/apps/YOUR_APP_ID/show
func SetConsumerSecret(consumer_secret string) {
	oauthClient.Credentials.Secret = consumer_secret
}

//AuthorizationURL generates the authorization URL for the first part of the OAuth handshake.
//Redirect the user to this URL.
//This assumes that the consumer key has already been set (using SetConsumerKey).
func AuthorizationURL(callback string) (string, *oauth.Credentials, error) {
	tempCred, err := oauthClient.RequestTemporaryCredentials(http.DefaultClient, callback, nil)
	if err != nil {
		return "", nil, err
	}
	return oauthClient.AuthorizationURL(tempCred, nil), tempCred, nil
}

func GetCredentials(tempCred *oauth.Credentials, verifier string) (*oauth.Credentials, url.Values, error) {
  return oauthClient.RequestToken(http.DefaultClient, tempCred, verifier)
}

func cleanValues(v url.Values) url.Values {
	if v == nil {
		return url.Values{}
	}
	return v
}

// apiGet issues a GET request to the Twitter API and decodes the response JSON to data.
func (c TwitterApi) apiGet(urlStr string, form url.Values, data interface{}) error {
	resp, err := oauthClient.Get(http.DefaultClient, c.Credentials, urlStr, form)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return decodeResponse(resp, data)
}

// apiPost issues a POST request to the Twitter API and decodes the response JSON to data.
func (c TwitterApi) apiPost(urlStr string, form url.Values, data interface{}) error {
	resp, err := oauthClient.Post(http.DefaultClient, c.Credentials, urlStr, form)
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
