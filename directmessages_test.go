package anaconda_test

import (
	"fmt"
	"github.com/jpibarra1130/anaconda"
	"net/url"
	"os"
	"testing"
)

var CONSUMER_KEY = os.Getenv("CONSUMER_KEY")
var CONSUMER_SECRET = os.Getenv("CONSUMER_SECRET")
var ACCESS_TOKEN = os.Getenv("ACCESS_TOKEN")
var ACCESS_TOKEN_SECRET = os.Getenv("ACCESS_TOKEN_SECRET")

var api *anaconda.TwitterApi

func init() {
	// Initialize api so it can be used even when invidual tests are run in isolation

	fmt.Println("Consumer key", CONSUMER_KEY)
	anaconda.SetConsumerKey(CONSUMER_KEY)
	anaconda.SetConsumerSecret(CONSUMER_SECRET)
	api = anaconda.NewTwitterApi(ACCESS_TOKEN, ACCESS_TOKEN_SECRET)
}

// Test that the client can be used to throttle to an arbitrary duration
func Test_TwitterApi_SendDirectMessage(t *testing.T) {
	userId := "paolo_ibarra"
	text := "This is a test message"
	v := url.Values{}
	v.Set("screen_name", userId)
	v.Set("text", text)

	directMessage, err := api.SendDirectMessage(v)

	if err != nil {
		panic(err)
	}

	fmt.Println(directMessage)
}
