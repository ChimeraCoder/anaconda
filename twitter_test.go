package anaconda

import (
    "github.com/ChimeraCoder/anaconda"
	"fmt"
)

func ExampleSearch() {

	anaconda.SetConsumerKey("your-consumer-key")
	anaconda.SetConsumerSecret("your-consumer-secret")
	api := anaconda.NewTwitterApi("your-access-token", "your-access-token-secret")
	search_result, err := api.GetSearch("golang", nil)
    if err != nil{
        panic(err)
    }
	for _, tweet := range search_result {
		fmt.Print(tweet.Text)
	}
}
