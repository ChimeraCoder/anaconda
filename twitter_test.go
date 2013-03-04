package twitter

import (
	"fmt"
)

func ExampleSearch() {

	twitter.SetConsumerKey("your-consumer-key")
	twitter.SetConsumerSecret("your-consumer-secret")
	api := twitter.NewTwitterApi("your-access-token", "your-access-token-secret")
	search_result, err := api.GetSearch("golang", url.Values{})
	for _, tweet := range search_result {
		fmt.Print(tweet.Text)
	}
}
