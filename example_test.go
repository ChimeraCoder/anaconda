package anaconda_test

import (
	"fmt"
	"time"

	"github.com/ChimeraCoder/anaconda"
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

func ExampleTwitterApi_PublicStreamFilter() {
	anaconda.SetConsumerKey("your-consumer-key")
	anaconda.SetConsumerSecret("your-consumer-secret")
	api := anaconda.NewTwitterApi("your-access-token", "your-access-token-secret")

	filter := map[string][]string{
		"track": []string{"#anaconda"},
		// Other options
		// "follow": []string{},
		// "locations": []string {},
	}
	tweetStream := api.PublicStreamFilter(filter)

	for t := range tweetStream.C {
		switch v := t.(type) {
		case anaconda.Tweet:
			fmt.Printf("%-15s: %s\n", v.User.ScreenName, v.Text)
		default:
			fmt.Printf("Got Something other than a tweet %s\n", v)
		}
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
