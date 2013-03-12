package anaconda

import (
	"encoding/json"
	"github.com/garyburd/go-oauth/oauth"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

var oauthClient = oauth.Client{
	TemporaryCredentialRequestURI: "http://api.twitter.com/oauth/request_token",
	ResourceOwnerAuthorizationURI: "http://api.twitter.com/oauth/authenticate",
	TokenRequestURI:               "http://api.twitter.com/oauth/access_token",
}

const (

	//As defined on https://dev.twitter.com/docs/error-codes-responses
	TwitterErrorDoesNotExist            = 34
	TwitterErrorRateLimitExceeded       = 88
	TwitterErrorInvalidToken            = 89
	TwitterErrorOverCapacity            = 130
	TwitterErrorInternalError           = 131
	TwitterErrorCouldNotAuthenticateYou = 135
	TwitterErrorBadAuthenticationData   = 215
)

type TwitterApi struct {
	Credentials *oauth.Credentials
}

type ApiError struct {
	errorString   string
	httpStatus    int
	TwitterErrors error  //If non-nil, this will be a TwitterError struct.
	requestUrl    string //If this was in response to a request, which endpoint?
}

func (e ApiError) Error() string {
	return e.errorString
}

func (e ApiError) HttpCode() int {
	return e.httpStatus
}

func (e ApiError) Code() int {
	return e.httpStatus
}

//TwitterError corresponds to the JSON errors that Twitter may return in API queries
type TwitterError struct {
	Message   string
	Code      int
	NextError error //Will be non-nil if Twitter returned more than one error
}

//OrMap returns true if the function evalutes to true on any TwitterError later in the list
func (c TwitterError) OrMap(f func(TwitterError) bool) bool {
	if f(c) {
		return true
	}
	if c.NextError == nil {
		return false
	}
	return c.NextError.(TwitterError).OrMap(f)
}

//ContainsError returns true if the current error or any later error in the list matches the error code specified.
func (err TwitterError) ContainsError(code int) bool {
	return err.OrMap(func(e TwitterError) bool {
		return e.Code == code
	})
}

type twitterErrorResponse struct {
	Errors []TwitterError
}

func (e TwitterError) Error() string {
	return e.Message
}

//NewTwitterApi takes an user-specific access token and secret and returns a TwitterApi struct for that user.
//The TwitterApi struct can be used for accessing any of the endpoints available.
func NewTwitterApi(access_token string, access_token_secret string) TwitterApi {
	return TwitterApi{&oauth.Credentials{Token: access_token, Secret: access_token_secret}}
}

//SetConsumerKey will set the application-specific consumer_key used in the initial OAuth process
//This key is listed on https://dev.twitter.com/apps/{APP_ID}/show
func SetConsumerKey(consumer_key string) {
	oauthClient.Credentials.Token = consumer_key
}

//SetConsumerSecret will set the application-specific secret used in the initial OAuth process
//This secret is listed on https://dev.twitter.com/apps/{APP_ID}/show
func SetConsumerSecret(consumer_secret string) {
	oauthClient.Credentials.Secret = consumer_secret
}

//AuthorizationURL generates the authorization URL for the first part of the OAuth handshake.
//Redirect the user to this URL.
//This assumes that the consumer key has already been set (using SetConsumerKey).
func AuthorizationURL(callback string) (string, error) {
	tempCred, err := oauthClient.RequestTemporaryCredentials(http.DefaultClient, callback, nil)
	if err != nil {
		return "", err
	}
	return oauthClient.AuthorizationURL(tempCred, nil), nil
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
func decodeResponse(resp *http.Response, data interface{}) (err error) {
	if resp.StatusCode != 200 {
		p, _ := ioutil.ReadAll(resp.Body)

		//Decode the error message(s) sent by Twitter
		var err_resp twitterErrorResponse
		if err := json.Unmarshal(p, &err_resp); err != nil {
			return err
		}

		for i := 0; i < (len(err_resp.Errors) - 1); i++ {
			err_resp.Errors[i].NextError = err_resp.Errors[i+1]
		}
		err = err_resp.Errors[0]
		log.Printf("We're passing in errors %+v", err)

		return ApiError{string(p), resp.StatusCode, err, resp.Request.URL.String()}
	}
	return json.NewDecoder(resp.Body).Decode(data)
}
