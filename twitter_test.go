package anaconda_test

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"os"
	"reflect"
	"testing"
	"time"
)

var CONSUMER_KEY = os.Getenv("CONSUMER_KEY")
var CONSUMER_SECRET = os.Getenv("CONSUMER_SECRET")
var ACCESS_TOKEN = os.Getenv("ACCESS_TOKEN")
var ACCESS_TOKEN_SECRET = os.Getenv("ACCESS_TOKEN_SECRET")

var api *anaconda.TwitterApi

func init() {
	// Initialize api so it can be used even when invidual tests are run in isolation
	anaconda.SetConsumerKey(CONSUMER_KEY)
	anaconda.SetConsumerSecret(CONSUMER_SECRET)
	api = anaconda.NewTwitterApi(ACCESS_TOKEN, ACCESS_TOKEN_SECRET)
}

// Test_TwitterCredentials tests that non-empty Twitter credentials are set
// Without this, all following tests will fail
func Test_TwitterCredentials(t *testing.T) {
	if CONSUMER_KEY == "" || CONSUMER_SECRET == "" || ACCESS_TOKEN == "" || ACCESS_TOKEN_SECRET == "" {
		t.Errorf("Credentials are invalid: at least one is empty")
	}
}

// Test that creating a TwitterApi client creates a client with non-empty OAuth credentials
func Test_TwitterApi_NewTwitterApi(t *testing.T) {
	anaconda.SetConsumerKey(CONSUMER_KEY)
	anaconda.SetConsumerSecret(CONSUMER_SECRET)
	api = anaconda.NewTwitterApi(ACCESS_TOKEN, ACCESS_TOKEN_SECRET)

	if api.Credentials == nil {
		t.Errorf("Twitter Api client has empty (nil) credentials")
	}
}

// Test that the GetSearch function actually works and returns non-empty results
func Test_TwitterApi_GetSearch(t *testing.T) {
	search_result, err := api.GetSearch("golang", nil)
	if err != nil {
		t.Errorf("GetSearch yielded error %s", err.Error())
		panic(err)
	}

	// Unless something is seriously wrong, there should be at least two tweets
	if len(search_result) < 2 {
		t.Errorf("Expected 2 or more tweets, and found %d", len(search_result))
	}

	// Check that at least one tweet is non-empty
	for _, tweet := range search_result {
		if tweet.Text != "" {
			return
		}
		fmt.Print(tweet.Text)
	}

	t.Errorf("All %d tweets had empty text", len(search_result))
}

// Test that setting the delay actually changes the stored delay value
func Test_TwitterApi_SetDelay(t *testing.T) {
	const OLD_DELAY = 1 * time.Second
	const NEW_DELAY = 20 * time.Second
	api.EnableRateLimiting(OLD_DELAY, 4)

	delay := api.GetDelay()
	if delay != OLD_DELAY {
		t.Errorf("Expected initial delay to be the default delay (%s)", anaconda.DEFAULT_DELAY.String())
	}

	api.SetDelay(NEW_DELAY)

	if newDelay := api.GetDelay(); newDelay != NEW_DELAY {
		t.Errorf("Attempted to set delay to %s, but delay is now %s (original delay: %s)", NEW_DELAY, newDelay, delay)
	}
}

// Test that the client can be used to throttle to an arbitrary duration
func Test_TwitterApi_Throttling(t *testing.T) {
	const MIN_DELAY = 30 * time.Second

	api.EnableRateLimiting(MIN_DELAY, 5)
	oldDelay := api.GetDelay()
	api.SetDelay(MIN_DELAY)

	now := time.Now()
	_, err := api.GetSearch("golang", nil)
	if err != nil {
		t.Errorf("GetSearch yielded error %s", err.Error())
	}
	_, err = api.GetSearch("anaconda", nil)
	if err != nil {
		t.Errorf("GetSearch yielded error %s", err.Error())
	}
	after := time.Now()

	if difference := after.Sub(now); difference < (30 * time.Second) {
		t.Errorf("Expected delay of at least %d. Actual delay: %s", MIN_DELAY.String(), difference.String())
	}

	// Reset the delay to its previous value
	api.SetDelay(oldDelay)
}

func Test_TwitterApi_TwitterErrorDoesNotExist(t *testing.T) {

	// Try fetching a tweet that no longer exists (was deleted)
	const DELETED_TWEET_ID = 404409873170841600

	tweet, err := api.GetTweet(DELETED_TWEET_ID, nil)
	if err == nil {
		t.Errorf("Expected an error when fetching tweet with id %d but got none - tweet object is %+v", DELETED_TWEET_ID, tweet)
	}

	apiErr, ok := err.(*anaconda.ApiError)
	if !ok {
		t.Errorf("Expected an *anaconda.ApiError, and received error message %s, (%+v)", err.Error(), err)
	}

	terr, ok := apiErr.Decoded.First().(anaconda.TwitterError)

	if !ok {
		t.Errorf("TwitterErrorResponse.First() should return value of type TwitterError, not %s", reflect.TypeOf(apiErr.Decoded.First()))
	}

	if code := terr.Code; code != anaconda.TwitterErrorDoesNotExist {
		if code == anaconda.TwitterErrorRateLimitExceeded {
			t.Errorf("Rate limit exceeded during testing - received error code %d instead of %d", anaconda.TwitterErrorRateLimitExceeded, anaconda.TwitterErrorDoesNotExist)
		} else {

			t.Errorf("Expected Twitter to return error code %d, and instead received error code %d", anaconda.TwitterErrorDoesNotExist, code)
		}
	}
}

// Test that a valid user can be fetched
// and that unmarshalling works properly
func Test_GetUser(t *testing.T) {
	const username = "chimeracoder"

	users, err := api.GetUsersLookup(username, nil)
	if err != nil {
		t.Errorf("GetUsersLookup returned error: %s", err.Error())
	}

	if len(users) != 1 {
		t.Errorf("Expected one user and received %d", len(users))
	}

	if !reflect.DeepEqual(users[0], anaconda.TwitterUser{}) {
		t.Errorf("Received %+v", users[0])
	}

}

// Test that a valid tweet can be fetched properly
// and that unmarshalling of tweet works without error
func Test_GetTweet(t *testing.T) {
	const tweetId = 303777106620452864
	const tweetText = `golang-syd is in session. Dave Symonds is now talking about API design and protobufs. #golang http://t.co/eSq3ROwu`

	tweet, err := api.GetTweet(tweetId, nil)
	if err != nil {
		t.Errorf("GetTweet returned error: %s", err.Error())
	}

	if tweet.Text != tweetText {
		t.Errorf("Tweet %d contained incorrect text. Received: %s", tweetId, tweetText)
	}
}
