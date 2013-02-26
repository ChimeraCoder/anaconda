package twitter

type DirectMessage struct {
	text     string
	entities struct {
		urls          []interface{}
		hashtags      []interface{}
		user_mentions []interface{}
	}
	id_str                string
	created_at            string
	recipient_id          float64
	id                    float64
	recipient_screen_name string
	sender                struct {
		name                               string
		default_profile_image              bool
		profile_image_url_https            string
		notifications                      bool
		protected                          bool
		id_str                             string
		profile_background_color           string
		created_at                         string
		default_profile                    bool
		url                                string
		time_zone                          string
		id                                 float64
		verified                           bool
		profile_link_color                 string
		profile_image_url                  string
		profile_use_background_image       bool
		favourites_count                   float64
		profile_background_image_url_https string
		profile_sidebar_fill_color         string
		utc_offset                         float64
		is_translator                      bool
		follow_request_sent                bool
		following                          bool
		profile_background_tile            bool
		show_all_inline_media              bool
		profile_text_color                 string
		lang                               string
		statuses_count                     float64
		contributors_enabled               bool
		friends_count                      float64
		geo_enabled                        bool
		description                        string
		profile_sidebar_border_color       string
		screen_name                        string
		listed_count                       float64
		followers_count                    float64
		location                           string
		profile_background_image_url       string
	}
	sender_id float64
	recipient struct {
		profile_text_color                 string
		lang                               string
		statuses_count                     float64
		contributors_enabled               bool
		friends_count                      float64
		geo_enabled                        bool
		description                        string
		profile_sidebar_border_color       string
		screen_name                        string
		listed_count                       float64
		followers_count                    float64
		location                           string
		profile_background_image_url       string
		name                               string
		default_profile_image              bool
		profile_image_url_https            string
		notifications                      bool
		protected                          bool
		id_str                             string
		profile_background_color           string
		created_at                         string
		default_profile                    bool
		url                                string
		time_zone                          string
		id                                 float64
		verified                           bool
		profile_link_color                 string
		profile_image_url                  string
		profile_use_background_image       bool
		favourites_count                   float64
		profile_background_image_url_https string
		profile_sidebar_fill_color         string
		utc_offset                         float64
		is_translator                      bool
		follow_request_sent                bool
		following                          bool
		profile_background_tile            bool
		show_all_inline_media              bool
	}
	sender_screen_name string
}
