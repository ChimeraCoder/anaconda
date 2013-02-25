package twitter

import (
    "net/url"
    "fmt"
    "strconv"
)

func (a TwitterApi) GetTweet(id int64, v url.Values) (tweet Tweet, err error){
    v.Set("id", strconv.FormatInt(id,10))
    err = a.apiGet("https://api.twitter.com/1.1/statuses/show.json", v, &tweet)
    return
}

func (a TwitterApi) GetRetweets(id int64, v url.Values) (tweets []Tweet, err error){
    err = a.apiGet(fmt.Sprintf("https://api.twitter.com/1.1/statuses/retweets/%d.json", id), v, &tweets)
    return
}




