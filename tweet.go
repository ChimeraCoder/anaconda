package anaconda

import (
	"time"
)

type Tweet struct {
	Source        string
	Id            int64
	Retweeted     bool
	Favorited     bool
	User          TwitterUser
	Truncated     bool
	Text          string
	Retweet_count int64
	Id_str        string
	Created_at    string
	Entities      TwitterEntities
}

// CreatedAtTime is a convenience wrapper that returns the Created_at time, parsed as a time.Time struct
func (t Tweet) CreatedAtTime() (time.Time, error) {
	return time.Parse(time.RubyDate, t.Created_at)
}
