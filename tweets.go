package twitter

import (
    "net/url"
    "strconv"
)

func (a TwitterApi) GetTweet(id int64, v url.Values) (tweet Tweet, err error){
    v.Set("id", strconv.FormatInt(id,10))
    err = a.apiGet("https://api.twitter.com/1.1/statuses/show.json", v, &tweet)
    return
}

