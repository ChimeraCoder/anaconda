package anaconda_test

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"os"
	"testing"
)

func Test_TwitterApi_GetSearch(t *testing.T) {
	anaconda.SetConsumerKey(os.Getenv("CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("CONSUMER_SECRET"))
	api := anaconda.NewTwitterApi(os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_TOKEN_SECRET"))
	search_result, err := api.GetSearch("golang", nil)
	if err != nil {
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
