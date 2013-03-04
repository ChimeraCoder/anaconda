package twitter

type Tweet struct {
	Source    string
	Id        int64
	Retweeted bool
	Favorited bool
	//User          TwitterUser
	Truncated     bool
	Text          string
	Retweet_count int64
	Id_str        string
	Created_at    string
	Entities      TwitterEntities
}
