package anaconda

type Retweet struct {
	Favorited     *bool
	User          TwitterUser
	Truncated     *bool
	Text          *string
	Retweet_count *float64
	Entities      struct {
		Urls          []interface{}
		Hashtags      []interface{}
		User_mentions []interface{}
	}
	Retweeted_status struct {
		Favorited     *bool
		User          TwitterUser
		Truncated     *bool
		Text          *string
		Retweet_count *float64
		Entities      struct {
			Hashtags      []interface{}
			User_mentions []interface{}
			Urls          []interface{}
		}
		Id_str     *string
		Created_at *string
		Source     *string
		Id         *float64
		Retweeted  *bool
	}
	Id_str     *string
	Created_at *string
	Source     *string
	Id         *float64
	Retweeted  *bool
}
