package anaconda

import "net/url"

type Media struct {
	MediaID       int64  `json:"media_id"`
	MediaIDString string `json:"media_id_string"`
	Size          int    `json:"size"`
	Image         Image  `json:"image"`
}

type Image struct {
	W         int    `json:"w"`
	H         int    `json:"h"`
	ImageType string `json:"image_type"`
}

func (a TwitterApi) UploadMedia(base64String string) (media Media, err error) {
	v := url.Values{}
	v.Set("media", base64String)

	var mediaResponse Media

	response_ch := make(chan response)
	a.queryQueue <- query{UploadBaseUrl + "/media/upload.json", v, &mediaResponse, _POST, response_ch}
	return mediaResponse, (<-response_ch).err
}
