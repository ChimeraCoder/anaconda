package anaconda

type UrlItem struct{
	Indices      []int
	Url          string
	Display_url  string
	Expanded_url string
}

type UrlEntity struct {
	Urls []UrlItem
}

type Entities struct {
	Hashtags []struct {
		Indices []int
		Text    string
	}
	Urls          []UrlItem
	Url           UrlEntity
	User_mentions []struct {
		Name        string
		Indices     []int
		Screen_name string
		Id          int64
		Id_str      string
	}
	Media []EntityMedia
	Description   UrlEntity
}

type EntityMedia struct {
	Id                   int64
	Id_str               string
	Media_url            string
	Media_url_https      string
	Url                  string
	Display_url          string
	Expanded_url         string
	Sizes                MediaSizes
	Source_status_id     int64
	Source_status_id_str string
	Type                 string
	Indices              []int
	VideoInfo            VideoInfo `json:"video_info"`
	ExtAltText           string    `json:"ext_alt_text"`
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
