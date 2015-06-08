package anaconda_test

import (
	"fmt"
	"net/url"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/ChimeraCoder/anaconda"
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
	if len(search_result.Statuses) < 2 {
		t.Errorf("Expected 2 or more tweets, and found %d", len(search_result.Statuses))
	}

	// Check that at least one tweet is non-empty
	for _, tweet := range search_result.Statuses {
		if tweet.Text != "" {
			return
		}
		fmt.Print(tweet.Text)
	}

	t.Errorf("All %d tweets had empty text", len(search_result.Statuses))
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

	// If all attributes are equal to the zero value for that type,
	// then the original value was not valid
	if reflect.DeepEqual(users[0], anaconda.User{}) {
		t.Errorf("Received %#v", users[0])
	}
}

func Test_GetFavorites(t *testing.T) {
	v := url.Values{}
	v.Set("screen_name", "chimeracoder")
	favorites, err := api.GetFavorites(v)
	if err != nil {
		t.Errorf("GetFavorites returned error: %s", err.Error())
	}

	if len(favorites) == 0 {
		t.Errorf("GetFavorites returned no favorites")
	}

	if reflect.DeepEqual(favorites[0], anaconda.Tweet{}) {
		t.Errorf("GetFavorites returned %d favorites and the first one was empty", len(favorites))
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

	// Check the entities
	expectedEntities := anaconda.Entities{Hashtags: []struct {
		Indices []int
		Text    string
	}{struct {
		Indices []int
		Text    string
	}{Indices: []int{86, 93}, Text: "golang"}}, Urls: []struct {
		Indices      []int
		Url          string
		Display_url  string
		Expanded_url string
	}{}, User_mentions: []struct {
		Name        string
		Indices     []int
		Screen_name string
		Id          int64
		Id_str      string
	}{}, Media: []struct {
		Id              int64
		Id_str          string
		Media_url       string
		Media_url_https string
		Url             string
		Display_url     string
		Expanded_url    string
		Sizes           anaconda.MediaSizes
		Type            string
		Indices         []int
	}{struct {
		Id              int64
		Id_str          string
		Media_url       string
		Media_url_https string
		Url             string
		Display_url     string
		Expanded_url    string
		Sizes           anaconda.MediaSizes
		Type            string
		Indices         []int
	}{Id: 303777106628841472, Id_str: "303777106628841472", Media_url: "http://pbs.twimg.com/media/BDc7q0OCEAAoe2C.jpg", Media_url_https: "https://pbs.twimg.com/media/BDc7q0OCEAAoe2C.jpg", Url: "http://t.co/eSq3ROwu", Display_url: "pic.twitter.com/eSq3ROwu", Expanded_url: "http://twitter.com/golang/status/303777106620452864/photo/1", Sizes: anaconda.MediaSizes{Medium: anaconda.MediaSize{W: 600, H: 450, Resize: "fit"}, Thumb: anaconda.MediaSize{W: 150, H: 150, Resize: "crop"}, Small: anaconda.MediaSize{W: 340, H: 255, Resize: "fit"}, Large: anaconda.MediaSize{W: 1024, H: 768, Resize: "fit"}}, Type: "photo", Indices: []int{94, 114}}}}
	if !reflect.DeepEqual(tweet.Entities, expectedEntities) {
		t.Errorf("Tweet entities differ")
	}

}

// This assumes that the current user has at least two pages' worth of followers
func Test_GetFollowersListAll(t *testing.T) {
	result := api.GetFollowersListAll(nil)
	i := 0

	for page := range result {
		if i == 2 {
			return
		}

		if page.Error != nil {
			t.Errorf("Receved error from GetFollowersListAll: %s", page.Error)
		}

		if page.Followers == nil || len(page.Followers) == 0 {
			t.Errorf("Received invalid value for page %d of followers: %v", i, page.Followers)
		}
		i++
	}
}

// This assumes that the current user has at least two pages' worth of friends
func Test_GetFriendsIdsAll(t *testing.T) {
	result := api.GetFriendsIdsAll(nil)
	i := 0

	for page := range result {
		if i == 2 {
			return
		}

		if page.Error != nil {
			t.Errorf("Receved error from GetFriendsIdsAll : %s", page.Error)
		}

		if page.Ids == nil || len(page.Ids) == 0 {
			t.Errorf("Received invalid value for page %d of friends : %v", i, page.Ids)
		}
		i++
	}
}


// Test that setting the delay actually changes the stored delay value
func Test_TwitterApi_SetDelay(t *testing.T) {
	const OLD_DELAY = 1 * time.Second
	const NEW_DELAY = 20 * time.Second
	api.EnableThrottling(OLD_DELAY, 4)

	delay := api.GetDelay()
	if delay != OLD_DELAY {
		t.Errorf("Expected initial delay to be the default delay (%s)", anaconda.DEFAULT_DELAY.String())
	}

	api.SetDelay(NEW_DELAY)

	if newDelay := api.GetDelay(); newDelay != NEW_DELAY {
		t.Errorf("Attempted to set delay to %s, but delay is now %s (original delay: %s)", NEW_DELAY, newDelay, delay)
	}
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

	if code := terr.Code; code != anaconda.TwitterErrorDoesNotExist && code != anaconda.TwitterErrorDoesNotExist2 {
		if code == anaconda.TwitterErrorRateLimitExceeded {
			t.Errorf("Rate limit exceeded during testing - received error code %d instead of %d", anaconda.TwitterErrorRateLimitExceeded, anaconda.TwitterErrorDoesNotExist)
		} else {

			t.Errorf("Expected Twitter to return error code %d, and instead received error code %d", anaconda.TwitterErrorDoesNotExist, code)
		}
	}
}

// Test that the client can be used to throttle to an arbitrary duration
func Test_TwitterApi_Throttling(t *testing.T) {
	const MIN_DELAY = 15 * time.Second

	api.EnableThrottling(MIN_DELAY, 5)
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

	if difference := after.Sub(now); difference < MIN_DELAY {
		t.Errorf("Expected delay of at least %s. Actual delay: %s", MIN_DELAY.String(), difference.String())
	}

	// Reset the delay to its previous value
	api.SetDelay(oldDelay)
}

func Test_DMScreenName(t *testing.T) {
	to, err := api.GetSelf(url.Values{})
	if err != nil {
		t.Error(err)
	}
	_, err = api.PostDMToScreenName("Test the anaconda lib", to.ScreenName)
	if err != nil {
		t.Error(err)
		return
	}
}
