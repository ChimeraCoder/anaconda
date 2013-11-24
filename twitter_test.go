package anaconda_test

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"log"
	"os"
	"testing"
	//"time"
)

var CONSUMER_KEY = os.Getenv("CONSUMER_KEY")
var CONSUMER_SECRET = os.Getenv("CONSUMER_SECRET")
var ACCESS_TOKEN = os.Getenv("ACCESS_TOKEN")
var ACCESS_TOKEN_SECRET = os.Getenv("ACCESS_TOKEN_SECRET")

var api *anaconda.TwitterApi

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

/*
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
	const NEW_DELAY = 20 * time.Second

	delay := api.GetDelay()
	if delay != anaconda.DEFAULT_DELAY {
		t.Errorf("Expected initial delay to be the default delay (%s)", anaconda.DEFAULT_DELAY.String())
	}

	api.SetDelay(NEW_DELAY)

	if newDelay := api.GetDelay(); newDelay != NEW_DELAY {
		t.Errorf("Attempted to set delay to %s, but delay is now %s (original delay: %s)", NEW_DELAY, newDelay, delay)
	}
}

// Test that the client can be used to throttle to an arbitrary duration
func Test_TwitterApi_Throttling(t *testing.T) {
	const MIN_DELAY_SECONDS = 30

    oldDelay := api.GetDelay()
	api.SetDelay(MIN_DELAY_SECONDS * time.Second)

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
		t.Errorf("Expected delay of at least %d seconds. Actual delay: %s", MIN_DELAY_SECONDS, difference.String())
	}

	// Reset the delay to its previous value
	api.SetDelay(oldDelay)
}
**/

func Test_TwitterApi_TwitterErrorDoesNotExist(t *testing.T) {

	// Try fetching a tweet that no longer exists (was deleted)
	const DELETED_TWEET_ID = 404409873170841600

	tweet, err := api.GetTweet(DELETED_TWEET_ID, nil)
	if err == nil {
		t.Errorf("Expected an error when fetching tweet with id %d but got none - tweet object is %+v", DELETED_TWEET_ID, tweet)
	}

	terr, ok := err.(anaconda.TwitterError)
	if !ok {
		log.Print(terr.Error())
		t.Errorf("Expected a TwitterError struct, and received error message %s, (%+v)", terr.Error(), terr)
	}
	if code := terr.Code; code != anaconda.TwitterErrorDoesNotExist {
		if code == anaconda.TwitterErrorRateLimitExceeded {
			t.Errorf("Rate limit exceeded during testing - received error code %d instead of %d", anaconda.TwitterErrorRateLimitExceeded, anaconda.TwitterErrorDoesNotExist)
		} else {

			t.Errorf("Expected Twitter to return error code %d, and instead received error code %d", anaconda.TwitterErrorDoesNotExist, code)
		}
	}
}

func ExampleTwitterApi_GetSearch() {

	anaconda.SetConsumerKey("your-consumer-key")
	anaconda.SetConsumerSecret("your-consumer-secret")
	api := anaconda.NewTwitterApi("your-access-token", "your-access-token-secret")
	search_result, err := api.GetSearch("golang", nil)
	if err != nil {
		panic(err)
	}
	for _, tweet := range search_result {
		fmt.Print(tweet.Text)
	}
}
