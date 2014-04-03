package anaconda

type UrlEntity struct {
	Urls []struct {
		Indices      []int
		Url          string
		Display_url  string
		Expanded_url string
	}
}

type Entities struct {
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
	Url           UrlEntity
	User_mentions []struct {
		Name        string
		Indices     []int
		Screen_name string
		Id          int64
		Id_str      string
	}
	Media []struct {
		Id              int64
		Id_str          string
		Media_url       string
		Media_url_https string
		Url             string
		Display_url     string
		Expanded_url    string
		Sizes           MediaSizes
		Type            string
		Indices         []int
	}
}

type MediaSizes struct {
	Medium MediaSize
	Thumb  MediaSize
	Small  MediaSize
	Large  MediaSize
}

type MediaSize struct {
	W      int
	H      int
	Resize string
}
