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
//  for _ , tweet := range searchResult.Statuses {
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
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/ChimeraCoder/tokenbucket"
	"github.com/garyburd/go-oauth/oauth"
	"sync"
)

const (
	_GET          = iota
	_POST         = iota
	BaseUrlV1     = "https://api.twitter.com/1"
	BaseUrl       = "https://api.twitter.com/1.1"
	UploadBaseUrl = "https://upload.twitter.com/1.1"
	tempCredsReqURI = "https://api.twitter.com/oauth/request_token"
	resourceOwnerAuthURI = "https://api.twitter.com/oauth/authenticate"
	tokenReqURI = "https://api.twitter.com/oauth/access_token"
)

type clientWithCreds struct {
	client oauth.Client
	accessToken oauth.Credentials
}

type clientCredsQueue struct {
	ch chan clientWithCreds
	n int
	mu sync.Mutex
}


func (q *clientCredsQueue) Add(c clientWithCreds) {
	q.mu.Lock()
	q.n += 1
	q.mu.Unlock()
	q.ch <- c
}

func (q *clientCredsQueue) Take() clientWithCreds {
	ret := <-q.ch
	q.mu.Lock()
	q.n -= 1
	q.mu.Unlock()
	return ret
}

func (q *clientCredsQueue) String() string {
	q.mu.Lock()
	s := fmt.Sprintf("clientCredsQueue: %d clients in the queue\n", q.n)
	q.mu.Unlock()
	return s
}

var clientQueue = &clientCredsQueue{ch:make(chan clientWithCreds, 15000)}

type TwitterCredentials struct {
	AccessToken string
	TokenSecret string
	ConsumerKey string
	ConsumerSecret string
}

func AddCredentials(creds ...TwitterCredentials) {
	for _, cred := range creds {
		client := oauth.Client{
			TemporaryCredentialRequestURI: tempCredsReqURI,
			ResourceOwnerAuthorizationURI: resourceOwnerAuthURI,
			TokenRequestURI: tokenReqURI,
			Credentials: oauth.Credentials{
				Token: cred.ConsumerKey,
				Secret: cred.ConsumerSecret,
			},
		}
		tok := oauth.Credentials{
			Token: cred.AccessToken,
			Secret: cred.TokenSecret,
		}

		c:= clientWithCreds{
			client: client,
			accessToken: tok,
		}

		clientQueue.Add(c)
	}
}

type TwitterApi struct {
	clientWithCreds
	queryQueue           chan query
	bucket               *tokenbucket.Bucket
	returnRateLimitError bool
	HttpClient           *http.Client

	// Currently used only for the streaming API
	// and for checking rate-limiting headers
	// Default logger is silent
	Log Logger
}

type query struct {
	url         string
	form        url.Values
	data        interface{}
	method      int
	response_ch chan response
}

type response struct {
	data interface{}
	err  error
}

const DEFAULT_DELAY = 0 * time.Second
const DEFAULT_CAPACITY = 5

//NewTwitterApi takes an user-specific access token and secret and returns a TwitterApi struct for that user.
//The TwitterApi struct can be used for accessing any of the endpoints available.
func NewTwitterApi() *TwitterApi {
	//TODO figure out how much to buffer this channel
	//A non-buffered channel will cause blocking when multiple queries are made at the same time
	queue := make(chan query)
	c := &TwitterApi{
		clientWithCreds:      clientQueue.Take(), // Can block if not enough clients
		queryQueue:           queue,
		bucket:               nil,
		returnRateLimitError: false,
		HttpClient:           http.DefaultClient,
		Log:                  silentLogger{},
	}
	go c.throttledQuery()
	return c
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

func cleanValues(v url.Values) url.Values {
	if v == nil {
		return url.Values{}
	}
	return v
}

// apiGet issues a GET request to the Twitter API and decodes the response JSON to data.
func (c TwitterApi) apiGet(urlStr string, form url.Values, data interface{}) error {
	resp, err := c.client.Get(c.HttpClient, &c.accessToken, urlStr, form)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return decodeResponse(resp, data)
}

// apiPost issues a POST request to the Twitter API and decodes the response JSON to data.
func (c TwitterApi) apiPost(urlStr string, form url.Values, data interface{}) error {
	resp, err := c.client.Post(c.HttpClient, &c.accessToken, urlStr, form)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return decodeResponse(resp, data)
}

// decodeResponse decodes the JSON response from the Twitter API.
func decodeResponse(resp *http.Response, data interface{}) error {
	if resp.StatusCode != 200 {
		return newApiError(resp)
	}
	return json.NewDecoder(resp.Body).Decode(data)
}

func NewApiError(resp *http.Response) *ApiError {
	body, _ := ioutil.ReadAll(resp.Body)

	return &ApiError{
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		Body:       string(body),
		URL:        resp.Request.URL,
	}
}

//query executes a query to the specified url, sending the values specified by form, and decodes the response JSON to data
//method can be either _GET or _POST
func (c TwitterApi) execQuery(urlStr string, form url.Values, data interface{}, method int) error {
	switch method {
	case _GET:
		return c.apiGet(urlStr, form, data)
	case _POST:
		return c.apiPost(urlStr, form, data)
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

		err := c.execQuery(url, form, data, method)

		// Check if Twitter returned a rate-limiting error
		if err != nil {
			if apiErr, ok := err.(*ApiError); ok {
				if isRateLimitError, nextWindow := apiErr.RateLimitCheck(); isRateLimitError && !c.returnRateLimitError {
					c.Log.Info("Error is rate limited")
					c.Log.Info(apiErr.Error())

					// If this is a rate-limiting error, re-add the job to the queue
					// TODO it really should preserve order
					go func() {
						c.queryQueue <- q
					}()

					delay := nextWindow.Sub(time.Now())
					go func(cl clientWithCreds) {
						c.Log.Infof("Waiting %v to put client back on the queue", delay)
						<-time.After(delay)
						c.Log.Info("Delay over")
						clientQueue.Add(cl)
					}(c.clientWithCreds)

					c.Log.Info("Attempting to retrieve new client")
					c.clientWithCreds = clientQueue.Take()

					// Drain the bucket (start over fresh)
					if c.bucket != nil {
						c.bucket.Drain()
					}

					continue
				} else if apiErr.InvalidToken() {
					// If token is invalid, re-add to queue
					go func() {
						c.queryQueue <- q
					}()
					c.Log.Info("Token is invalid, discarding")

					c.Log.Info("Attempting to retrieve new client")
					c.clientWithCreds = clientQueue.Take()
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
}
