package anaconda

type UrlEntity struct {
	Urls []struct {
		Indices      []int  `json:"indices"`
		Url          string `json:"url"`
		Display_url  string `json:"display_url"`
		Expanded_url string `json:"expanded_url"`
	} `json:"urls"`
}

type Entities struct {
	Urls []struct {
		Indices      []int  `json:"indices"`
		Url          string `json:"url"`
		Display_url  string `json:"display_url"`
		Expanded_url string `json:"expanded_url"`
	} `json:"urls"`
	Hashtags []struct {
		Indices []int  `json:"indices"`
		Text    string `json:"text"`
	} `json:"hashtags"`
	Url           UrlEntity `json:"url"`
	User_mentions []struct {
		Name        string `json:"name"`
		Indices     []int  `json:"indices"`
		Screen_name string `json:"screen_name"`
		Id          int64  `json:"id"`
		Id_str      string `json:"id_str"`
	} `json:"user_mentions"`
	Media []EntityMedia `json:"media"`
}

type EntityMedia struct {
	Id                   int64      `json:"id"`
	Id_str               string     `json:"id_str"`
	Media_url            string     `json:"media_url"`
	Media_url_https      string     `json:"media_url_https"`
	Url                  string     `json:"url"`
	Display_url          string     `json:"display_url"`
	Expanded_url         string     `json:"expanded_url"`
	Sizes                MediaSizes `json:"sizes"`
	Source_status_id     int64      `json:"source_status_id"`
	Source_status_id_str string     `json:"source_status_id_str"`
	Type                 string     `json:"type"`
	Indices              []int      `json:"indices"`
	VideoInfo            VideoInfo  `json:"video_info"`
	ExtAltText           string     `json:"ext_alt_text"`
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

type Category struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
	Size int    `json:"size"`
}

type Suggestions struct {
	Category
	Users []User
}
