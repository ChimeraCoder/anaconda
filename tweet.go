package anaconda

import (
	"time"
)

type Tweet struct {
	Contributors []User `json:"contributors"`
	Coordinates  struct {
		Coordinates [2]float64 `json:"coordinates"`
		Type        string     `json:"type"`
	} `json:"coordinates"`
	CreatedAt            string      `json:"created_at"`
	Entities             Entities    `json:"entities"`
	FavoriteCount        int         `json:"favorite_count"`
	Favorited            bool        `json:"favorited"`
	FilterLevel          string      `json:"filter_level"`
	Id                   int64       `json:"id"`
	IdStr                string      `json:"id_str"`
	InReplyToScreenName  string      `json:"in_reply_to_screen_name"`
	InReplyToStatusID    int64       `json:"in_reply_to_status_id"`
	InReplyToStatusIdStr string      `json:"in_reply_to_status_id_str"`
	InReplyToUserID      int64       `json:"in_reply_to_user_id"`
	InReplyToUserIdStr   string      `json:"in_reply_to_user_id_str"`
	Lang                 string      `json:"lang"`
	Place                Place       `json:"place"`
	PossiblySensitive    bool        `json:"possibly_sensitive"`
	RetweetCount         int         `json:"retweet_count"`
	Retweeted            bool        `json:"retweeted"`
	RetweetedStatus      *Tweet      `json:"retweeted_status"`
	Source               string      `json:"source"`
	Scopes               interface{} `json:"scopes"`
	Text                 string      `json:"text"`
	Truncated            bool        `json:"truncated"`
	User                 User        `json:"user"`
	WithheldCopyright    bool        `json:"withheld_copyright"`
	WithheldInCountries  []string    `json:"withheld_in_countries"`
	WithheldScope        string      `json:"withheld_scope"`
	//Geo is deprecated, discourage usage
	//Geo                  interface{} `json:"geo"`
}

// CreatedAtTime is a convenience wrapper that returns the Created_at time, parsed as a time.Time struct
func (t Tweet) CreatedAtTime() (time.Time, error) {
	return time.Parse(time.RubyDate, t.CreatedAt)
}

func (t Tweet) HasCoordinates() bool {
	if t.Coordinates.Type == "Point" {
		return true
	}
	return false
}

// Latitude is a convenience wrapper that returns the latitude easily, not sure the best default return vaules
func (t Tweet) Latitude() float64 {
	if t.HasCoordinates() {
		return t.Coordinates.Coordinates[1]
	}
	return -9999
}

// Longitude is a convenience wrapper that returns the longitude easily, not sure the best default return values
func (t Tweet) Longitude() float64 {
	if t.HasCoordinates() {
		return t.Coordinates.Coordinates[0]
	}
	return -9999
}
