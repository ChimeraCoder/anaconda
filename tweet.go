package anaconda

import (
	"time"
)

type Tweet struct {
	Contributors         []int64     `json:"contributors"`
	Coordinates          interface{} `json:"coordinates"`
	CreatedAt            string      `json:"created_at"`
	Entities             Entities    `json:"entities"`
	FavoriteCount        int         `json:"favorite_count"`
	Favorited            bool        `json:"favorited"`
	Geo                  interface{} `json:"geo"`
	Id                   int64       `json:"id"`
	IdStr                string      `json:"id_str"`
	InReplyToScreenName  string      `json:"in_reply_to_screen_name"`
	InReplyToStatusID    int64       `json:"in_reply_to_status_id"`
	InReplyToStatusIdStr string      `json:"in_reply_to_status_id_str"`
	InReplyToUserID      int64       `json:"in_reply_to_user_id"`
	InReplyToUserIdStr   string      `json:"in_reply_to_user_id_str"`
	Place                Place       `json:"place"`
	PossiblySensitive    bool        `json:"possibly_sensitive"`
	RetweetCount         int         `json:"retweet_count"`
	Retweeted            bool        `json:"retweeted"`
	RetweetedStatus      *Tweet      `json:"retweeted_status"`
	Source               string      `json:"source"`
	Text                 string      `json:"text"`
	Truncated            bool        `json:"truncated"`
	User                 User        `json:"user"`
}

type ById []Tweet
// Implement sort interface
func (a ById) Len() int           { return len(a) }
func (a ById) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ById) Less(i, j int) bool { return a[i].Id < a[j].Id }

// CreatedAtTime is a convenience wrapper that returns the Created_at time, parsed as a time.Time struct
func (t Tweet) CreatedAtTime() (time.Time, error) {
	return time.Parse(time.RubyDate, t.CreatedAt)
}
