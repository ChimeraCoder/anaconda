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
	Entities      struct {
		Hashtags []struct {
			Indices []int
			Text    string
		}
		Urls []struct {
			Indices      []int
			Url          string
			Display_url  string
			Expanded_url string
		}
		User_mentions []struct {
			Name        string
			Indices     []int
			Screen_name string
			Id          int64
			Id_str      string
		}
	}
}
