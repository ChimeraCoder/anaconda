package anaconda_test

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"testing"
)

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

func ExampleErrorHandling() {
	api := anaconda.NewTwitterApi("invalid", "token")
	_, err := api.GetSearch("golang", nil)
	if err != nil {
		twitter_error := err.(anaconda.ApiError).TwitterErrors.(anaconda.TwitterError)
		if twitter_error.ContainsError(anaconda.TwitterErrorBadAuthenticationData) {
			//Some logic that fixes the authentication data
		}
	}
}

func TestTwitterErrorBadAuthenticationData(t *testing.T) {
	api := anaconda.NewTwitterApi("invalid", "token")
	_, err := api.GetSearch("golang", nil)
	if err != nil {
		aerr := err.(anaconda.ApiError)
		desired_error := anaconda.TwitterErrorBadAuthenticationData
		if !aerr.TwitterErrors.(anaconda.TwitterError).ContainsError(desired_error) {
			t.Errorf("Expected error code %d and received %+v", desired_error, aerr)
		}
	}

}
