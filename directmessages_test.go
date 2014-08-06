package anaconda_test

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
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

func Test_TwitterApi_SendDirectMessageScreenName(t *testing.T) {
	screenName := "paolo_ibarra"
	text := "This is a test message via screen name"

	directMessage, err := api.SendDirectMessageScreenName(screenName, text)

	if err != nil {
		panic(err)
	}

	fmt.Println(directMessage)
}

func Test_TwitterApi_SendDirectMessageUserId(t *testing.T) {
	userId := int64(186643967)
	text := "This is a test message via user id"

	directMessage, err := api.SendDirectMessageUserId(userId, text)

	if err != nil {
		panic(err)
	}

	fmt.Println(directMessage)
}
