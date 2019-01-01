package anaconda

type URLEntity struct {
	URLs []struct {
		Indices     []int  `json:"indices"`
		URL         string `json:"url"`
		DisplayURL  string `json:"display_url"`
		ExpandedURL string `json:"expanded_url"`
	} `json:"urls"`
}

type Entities struct {
	URLs []struct {
		Indices     []int  `json:"indices"`
		URL         string `json:"url"`
		DisplayURL  string `json:"display_url"`
		ExpandedURL string `json:"expanded_url"`
	} `json:"urls"`
	Hashtags []struct {
		Indices []int  `json:"indices"`
		Text    string `json:"text"`
	} `json:"hashtags"`
	URL          URLEntity `json:"url"`
	UserMentions []struct {
		Name       string `json:"name"`
		Indices    []int  `json:"indices"`
		ScreenName string `json:"screen_name"`
		ID         int64  `json:"id"`
		IDStr      string `json:"id_str"`
	} `json:"user_mentions"`
	Media []EntityMedia `json:"media"`
}

type EntityMedia struct {
	ID                int64      `json:"id"`
	IDStr             string     `json:"id_str"`
	MediaURL          string     `json:"media_url"`
	MediaURLHTTPS     string     `json:"media_url_https"`
	URL               string     `json:"url"`
	DisplayURL        string     `json:"display_url"`
	ExpandedURL       string     `json:"expanded_url"`
	Sizes             MediaSizes `json:"sizes"`
	SourceStatusID    int64      `json:"source_status_id"`
	SourceStatusIDStr string     `json:"source_status_id_str"`
	Type              string     `json:"type"`
	Indices           []int      `json:"indices"`
	VideoInfo         VideoInfo  `json:"video_info"`
	ExtAltText        string     `json:"ext_alt_text"`
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
	URL         string `json:"url"`
}

type Category struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
	Size int    `json:"size"`
}

type Suggestions struct {
	Category
	Users []User
}
