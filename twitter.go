// Package anaconda provides structs and functions for accessing version 1.1
// of the Twitter API.
//
// Successful API queries return native Go structs that can be used immediately,
// with no need for type assertions.
//
// # Authentication
//
// If you already have the access token (and secret) for your user (Twitter provides this for your own account on the developer portal), creating the client is simple:
//
//	anaconda.SetConsumerKey("your-consumer-key")
//	anaconda.SetConsumerSecret("your-consumer-secret")
//	api := anaconda.NewTwitterApi("your-access-token", "your-access-token-secret")
//
// # Queries
//
// Executing queries on an authenticated TwitterApi struct is simple.
//
//	searchResult, _ := api.GetSearch("golang", nil)
//	for _ , tweet := range searchResult.Statuses {
//	    fmt.Print(tweet.Text)
//	}
//
// Certain endpoints allow separate optional parameter; if desired, these can be passed as the final parameter.
//
//	v := url.Values{}
//	v.Set("count", "30")
//	result, err := api.GetSearch("golang", v)
//
// # Endpoints
//
// Anaconda implements most of the endpoints defined in the Twitter API documentation: https://dev.twitter.com/docs/api/1.1.
// For clarity, in most cases, the function name is simply the name of the HTTP method and the endpoint (e.g., the endpoint `GET /friendships/incoming` is provided by the function `GetFriendshipsIncoming`).
//
// In a few cases, a shortened form has been chosen to make life easier (for example, retweeting is simply the function `Retweet`)
//
// More detailed information about the behavior of each particular endpoint can be found at the official Twitter API documentation.
package anaconda

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ChimeraCoder/tokenbucket"
	"github.com/garyburd/go-oauth/oauth"
)

const (
	_GET          = iota
	_POST         = iota
	_DELETE       = iota
	_PUT          = iota
	_POSTBODY     = iota
	ClientTimeout = 20
	BaseUrlV1     = "https://api.twitter.com/1"
	BaseUrl       = "https://api.twitter.com/1.1"
	BaseUrlV2     = "https://api.twitter.com/2"
	UploadBaseUrl = "https://upload.twitter.com/1.1"
)

var (
	oauthCredentials oauth.Credentials
)

type TwitterApi struct {
	oauthClient          oauth.Client
	Credentials          *oauth.Credentials
	queryQueue           chan query
	queryQueueBody       chan queryBody
	bucket               *tokenbucket.Bucket
	returnRateLimitError bool
	HttpClient           *http.Client

	// Currently used only for the streaming API
	// and for checking rate-limiting headers
	// Default logger is silent
	Log Logger

	// used for testing
	// defaults to BaseUrl
	baseUrl   string
	baseUrlV2 string
}

type query struct {
	url         string
	form        url.Values
	data        interface{}
	method      int
	response_ch chan response
}
type queryBody struct {
	url         string
	form        url.Values
	data        interface{}
	method      int
	response_ch chan response
	Body        []byte
}

type response struct {
	data interface{}
	err  error
}

const DEFAULT_DELAY = 0 * time.Second
const DEFAULT_CAPACITY = 5

// NewTwitterApi takes an user-specific access token and secret and returns a TwitterApi struct for that user.
// The TwitterApi struct can be used for accessing any of the endpoints available.
func NewTwitterApi(access_token string, access_token_secret string) *TwitterApi {
	//TODO figure out how much to buffer this channel
	//A non-buffered channel will cause blocking when multiple queries are made at the same time
	_queue := make(chan query)
	_queryBody := make(chan queryBody)
	c := &TwitterApi{
		oauthClient: oauth.Client{
			TemporaryCredentialRequestURI: "https://api.twitter.com/oauth/request_token",
			ResourceOwnerAuthorizationURI: "https://api.twitter.com/oauth/authenticate",
			TokenRequestURI:               "https://api.twitter.com/oauth/access_token",
			Credentials:                   oauthCredentials,
		},
		Credentials: &oauth.Credentials{
			Token:  access_token,
			Secret: access_token_secret,
		},
		queryQueue:           _queue,
		queryQueueBody:       _queryBody,
		bucket:               nil,
		returnRateLimitError: false,
		HttpClient:           http.DefaultClient,
		Log:                  silentLogger{},
		baseUrl:              BaseUrl,
		baseUrlV2:            BaseUrlV2,
	}
	//Configure a timeout to HTTP client (DefaultClient has no default timeout, which may deadlock Mutex-wrapped uses of the lib.)
	c.HttpClient.Timeout = time.Duration(ClientTimeout * time.Second)
	go c.throttledQuery()
	go c.throttledQueryBody()
	return c
}

// NewTwitterApiWithCredentials takes an app-specific consumer key and secret, along with a user-specific access token and secret and returns a TwitterApi struct for that user.
// The TwitterApi struct can be used for accessing any of the endpoints available.
func NewTwitterApiWithCredentials(access_token string, access_token_secret string, consumer_key string, consumer_secret string) *TwitterApi {
	api := NewTwitterApi(access_token, access_token_secret)
	api.oauthClient.Credentials.Token = consumer_key
	api.oauthClient.Credentials.Secret = consumer_secret
	return api
}

// SetConsumerKey will set the application-specific consumer_key used in the initial OAuth process
// This key is listed on https://dev.twitter.com/apps/YOUR_APP_ID/show
func SetConsumerKey(consumer_key string) {
	oauthCredentials.Token = consumer_key
}

// SetConsumerSecret will set the application-specific secret used in the initial OAuth process
// This secret is listed on https://dev.twitter.com/apps/YOUR_APP_ID/show
func SetConsumerSecret(consumer_secret string) {
	oauthCredentials.Secret = consumer_secret
}

// ReturnRateLimitError specifies behavior when the Twitter API returns a rate-limit error.
// If set to true, the query will fail and return the error instead of automatically queuing and
// retrying the query when the rate limit expires
func (c *TwitterApi) ReturnRateLimitError(b bool) {
	c.returnRateLimitError = b
}

// Enable query throttling using the tokenbucket algorithm
func (c *TwitterApi) EnableThrottling(rate time.Duration, bufferSize int64) {
	c.bucket = tokenbucket.NewBucket(rate, bufferSize)
}

// Disable query throttling
func (c *TwitterApi) DisableThrottling() {
	c.bucket = nil
}

// SetDelay will set the delay between throttled queries
// To turn of throttling, set it to 0 seconds
func (c *TwitterApi) SetDelay(t time.Duration) {
	c.bucket.SetRate(t)
}

func (c *TwitterApi) GetDelay() time.Duration {
	return c.bucket.GetRate()
}

// SetBaseUrl is experimental and may be removed in future releases.
func (c *TwitterApi) SetBaseUrl(baseUrl string) {
	c.baseUrl = baseUrl
}

// AuthorizationURL generates the authorization URL for the first part of the OAuth handshake.
// Redirect the user to this URL.
// This assumes that the consumer key has already been set (using SetConsumerKey or NewTwitterApiWithCredentials).
func (c *TwitterApi) AuthorizationURL(callback string) (string, *oauth.Credentials, error) {
	tempCred, err := c.oauthClient.RequestTemporaryCredentials(http.DefaultClient, callback, nil)
	if err != nil {
		return "", nil, err
	}
	return c.oauthClient.AuthorizationURL(tempCred, nil), tempCred, nil
}

// GetCredentials gets the access token using the verifier received with the callback URL and the
// credentials in the first part of the handshake. GetCredentials implements the third part of the OAuth handshake.
// The returned url.Values holds the access_token, the access_token_secret, the user_id and the screen_name.
func (c *TwitterApi) GetCredentials(tempCred *oauth.Credentials, verifier string) (*oauth.Credentials, url.Values, error) {
	return c.oauthClient.RequestToken(http.DefaultClient, tempCred, verifier)
}

func defaultValues(v url.Values) url.Values {
	if v == nil {
		v = url.Values{}
	}
	v.Set("tweet_mode", "extended")
	return v
}

func cleanValues(v url.Values) url.Values {
	if v == nil {
		return url.Values{}
	}
	return v
}

// apiGet issues a GET request to the Twitter API and decodes the response JSON to data.
func (c TwitterApi) apiGet(urlStr string, form url.Values, data interface{}) error {
	form = defaultValues(form)
	resp, err := c.oauthClient.Get(c.HttpClient, c.Credentials, urlStr, form)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return decodeResponse(resp, data)
}

// apiPost issues a POST request to the Twitter API and decodes the response JSON to data.
func (c TwitterApi) apiPost(urlStr string, form url.Values, data interface{}) error {
	resp, err := c.oauthClient.Post(c.HttpClient, c.Credentials, urlStr, form)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return decodeResponse(resp, data)
}

// apiPost issues a POST request to the Twitter API and decodes the response JSON to data.
func (c TwitterApi) apiPostBody(urlStr string, form url.Values, body []byte, data interface{}) error {
	// Create a new HTTP request
	req, err := http.NewRequest("POST", urlStr, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	// Add any form values you want to the request
	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = form.Encode()

	// Add OAuth headers to the request
	c.oauthClient.SetAuthorizationHeader(req.Header, c.Credentials, "POST", req.URL, form)
	authorization := c.oauthClient.Header.Clone().Get("Authorization")
	req.Header.Add("Authorization", authorization)

	// Perform the request using the provided HTTP client
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode the response
	return decodeResponse(resp, data)
}

// apiDel issues a DELETE request to the Twitter API and decodes the response JSON to data.
func (c TwitterApi) apiDel(urlStr string, form url.Values, data interface{}) error {
	resp, err := c.oauthClient.Delete(c.HttpClient, c.Credentials, urlStr, form)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return decodeResponse(resp, data)
}

// apiPut issues a PUT request to the Twitter API and decodes the response JSON to data.
func (c TwitterApi) apiPut(urlStr string, form url.Values, data interface{}) error {
	resp, err := c.oauthClient.Put(c.HttpClient, c.Credentials, urlStr, form)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return decodeResponse(resp, data)
}

// decodeResponse decodes the JSON response from the Twitter API.
func decodeResponse(resp *http.Response, data interface{}) error {
	// Prevent memory leak in the case where the Response.Body is not used.
	// As per the net/http package, Response.Body still needs to be closed.
	defer resp.Body.Close()

	// Twitter returns deflate data despite the client only requesting gzip
	// data.  net/http automatically handles the latter but not the former:
	// https://github.com/golang/go/issues/18779
	if resp.Header.Get("Content-Encoding") == "deflate" {
		var err error
		resp.Body, err = zlib.NewReader(resp.Body)
		if err != nil {
			return err
		}
	}

	// according to dev.twitter.com, chunked upload append returns HTTP 2XX
	// so we need a special case when decoding the response
	if strings.HasSuffix(resp.Request.URL.String(), "upload.json") ||
		strings.Contains(resp.Request.URL.String(), "webhooks") {
		if resp.StatusCode == 204 {
			// empty response, don't decode
			return nil
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return newApiError(resp)
		}
	} else if resp.StatusCode != 200 {
		return newApiError(resp)
	}
	return json.NewDecoder(resp.Body).Decode(data)
}

func NewApiError(resp *http.Response) *ApiError {
	body, _ := io.ReadAll(resp.Body)

	return &ApiError{
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		Body:       string(body),
		URL:        resp.Request.URL,
	}
}

// query executes a query to the specified url, sending the values specified by form, and decodes the response JSON to data
// method can be either _GET or _POST
func (c TwitterApi) execQuery(urlStr string, form url.Values, data interface{}, method int, body []byte) error {
	switch method {
	case _GET:
		return c.apiGet(urlStr, form, data)
	case _POST:
		return c.apiPost(urlStr, form, data)
	case _DELETE:
		return c.apiDel(urlStr, form, data)
	case _PUT:
		return c.apiPut(urlStr, form, data)
	case _POSTBODY:
		return c.apiPostBody(urlStr, form, body, data)
	default:
		return fmt.Errorf("HTTP method not yet supported")
	}
}

// throttledQuery executes queries and automatically throttles them according to SECONDS_PER_QUERY
// It is the only function that reads from the queryQueue for a particular *TwitterApi struct
func (c *TwitterApi) throttledQuery() {
	for q := range c.queryQueue {
		url := q.url
		form := q.form
		data := q.data //This is where the actual response will be written
		method := q.method

		response_ch := q.response_ch

		if c.bucket != nil {
			<-c.bucket.SpendToken(1)
		}

		err := c.execQuery(url, form, data, method, nil)

		// Check if Twitter returned a rate-limiting error
		if err != nil {
			if apiErr, ok := err.(*ApiError); ok {
				if isRateLimitError, nextWindow := apiErr.RateLimitCheck(); isRateLimitError && !c.returnRateLimitError {
					c.Log.Info(apiErr.Error())

					// If this is a rate-limiting error, re-add the job to the queue
					// TODO it really should preserve order
					go func(q query) {
						c.queryQueue <- q
					}(q)

					delay := nextWindow.Sub(time.Now())
					<-time.After(delay)

					// Drain the bucket (start over fresh)
					if c.bucket != nil {
						c.bucket.Drain()
					}

					continue
				}
			}
		}

		response_ch <- response{data, err}
	}
}

func (c *TwitterApi) throttledQueryBody() {
	for q := range c.queryQueueBody {
		url := q.url
		form := q.form
		data := q.data //This is where the actual response will be written
		method := q.method
		_queryBody := q.Body

		response_ch := q.response_ch

		if c.bucket != nil {
			<-c.bucket.SpendToken(1)
		}

		err := c.execQuery(url, form, data, method, _queryBody)

		// Check if Twitter returned a rate-limiting error
		if err != nil {
			if apiErr, ok := err.(*ApiError); ok {
				if isRateLimitError, nextWindow := apiErr.RateLimitCheck(); isRateLimitError && !c.returnRateLimitError {
					c.Log.Info(apiErr.Error())

					// If this is a rate-limiting error, re-add the job to the queue
					// TODO it really should preserve order
					go func(q queryBody) {
						c.queryQueueBody <- q
					}(q)

					delay := nextWindow.Sub(time.Now())
					<-time.After(delay)

					// Drain the bucket (start over fresh)
					if c.bucket != nil {
						c.bucket.Drain()
					}

					continue
				}
			}
		}

		response_ch <- response{data, err}
	}
}

// Close query queue
func (c *TwitterApi) Close() {
	close(c.queryQueue)
	close(c.queryQueueBody)
}
