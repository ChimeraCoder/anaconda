package anaconda_test

import (
	"fmt"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"net/url"
)

// Initialize an client library for a given user.
// This only needs to be done *once* per user
func ExampleTwitterApi_InitializeClient() {
	api := anaconda.NewTwitterApiWithCredentials(ACCESS_TOKEN, ACCESS_TOKEN_SECRET, "your-consumer-key", "your-consumer-secret")
	fmt.Println(*api.Credentials)
}

func ExampleTwitterApi_GetSearch() {
	anaconda.SetConsumerKey("your-consumer-key")
	anaconda.SetConsumerSecret("your-consumer-secret")
	api := anaconda.NewTwitterApi("your-access-token", "your-access-token-secret")
	search_result, err := api.GetSearch("golang", nil)
	if err != nil {
		panic(err)
	}
	for _, tweet := range search_result.Statuses {
		fmt.Print(tweet.Text)
	}
}

// Throttling queries can easily be handled in the background, automatically
func ExampleTwitterApi_Throttling() {
	api := anaconda.NewTwitterApi("your-access-token", "your-access-token-secret")
	api.EnableThrottling(10*time.Second, 5)

	// These queries will execute in order
	// with appropriate delays inserted only if necessary
	golangTweets, err := api.GetSearch("golang", nil)
	anacondaTweets, err2 := api.GetSearch("anaconda", nil)

	if err != nil {
		panic(err)
	}
	if err2 != nil {
		panic(err)
	}

	fmt.Println(golangTweets)
	fmt.Println(anacondaTweets)
}

// Fetch a list of all followers without any need for managing cursors
// (Each page is automatically fetched when the previous one is read)
func ExampleTwitterApi_GetFollowersListAll() {
	pages := api.GetFollowersListAll(nil)
	for page := range pages {
		//Print the current page of followers
		fmt.Println(page.Followers)
	}
}

func ExampleTwitterApi_GetDMEventList() {
	anaconda.SetConsumerKey("your-consumer-key")
	anaconda.SetConsumerSecret("your-consumer-secret")
	api := anaconda.NewTwitterApi("your-access-token", "your-access-token-secret")

	v := url.Values{}
	v.Set("count", "50")
	v.Set("cursor", "next-cursor")
	result, err := api.GetDMEventList(v)
	if err != nil {
		panic(err)
	}

	fmt.Println(result.NextCursor)
	for _, event := range result.DMEvents {
		fmt.Println(event.Id)
		fmt.Println(event.MessageCreate.MessageData.Text)
	}
}

func ExampleTwitterApi_GetDMEventShow() {
	anaconda.SetConsumerKey("your-consumer-key")
	anaconda.SetConsumerSecret("your-consumer-secret")
	api := anaconda.NewTwitterApi("your-access-token", "your-access-token-secret")

	result, err := api.GetDMEventShow("your-event-id")
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}
