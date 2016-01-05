package anaconda

type Hashtags struct {
	Indices []int  `json:"indices"`
	Text    string `json:"text"`
}

type Urls struct {
	Indices     []int  `json:"indices"`
	Url         string `json:"url"`
	DisplayUrl  string `json:"display_url"`
	ExpandedUrl string `json:"expanded_url"`
}

type UserMentions struct {
	Name       string `json:"name"`
	Indices    []int  `json:"indices"`
	ScreenName string `json:"screen_name"`
	Id         int64  `json:"id"`
	IdStr      string `json:"id_str"`
}

type Entities struct {
	Hashtags     []Hashtags     `json:"hashtags"`
	Urls         []Urls         `json:"urls"`
	Media        []EntityMedia  `json:"media"`
	UserMentions []UserMentions `json:"user_mentions"`
}

type EntityMedia struct {
	Id                int64      `json:"id"`
	IdStr             string     `json:"id_str"`
	MediaUrl          string     `json:"media_url"`
	MediaUrlHttps     string     `json:"media_url_https"`
	Url               string     `json:"url"`
	DisplayUrl        string     `json:"display_url"`
	ExpandedUrl       string     `json:"expanded_url"`
	Sizes             MediaSizes `json:"sizes"`
	SourceStatusId    int64      `json:"source_status_id"`
	SourceStatusIdStr string     `json:"source_status_id_str"`
	Type              string     `json:"type"`
	Indices           []int      `json:"indices"`
	VideoInfo         VideoInfo  `json:"video_info"`
}

type MediaSizes struct {
	Medium MediaSize `json:"medium"`
	Thumb  MediaSize `json:"thumb"`
	Small  MediaSize `json:"small"`
	Large  MediaSize `json:"large"`
}

type MediaSize struct {
	W      int    `json:"w"`
	H      int    `json:"h"`
	Resize string `json:"resize"`
}

type VideoInfo struct {
	AspectRatio    []int     `json:"aspect_ratio"`
	DurationMillis int64     `json:"duration_millis"`
	Variants       []Variant `json:"variants"`
}

type Variant struct {
	Bitrate     int    `json:"bitrate"`
	ContentType string `json:"content_type"`
	Url         string `json:"url"`
}
