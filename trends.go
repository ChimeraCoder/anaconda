package anaconda

import "net/url"

type Location struct {
	Name  string `json:"name"`
	Woeid int    `json:"woeid"`
}

type Trend struct {
	Name            string `json:"name"`
	Query           string `json:"query"`
	Url             string `json:"url"`
	PromotedContent string `json:"promoted_content"`
}

type TrendResponse struct {
	Trends    []Trend    `json:"trends"`
	AsOf      string     `json:"as_of"`
	CreatedAt string     `json:"created_at"`
	Locations []Location `json:"locations"`
}

type MegaTrendResponse struct {
	Response TrendResponse
}

func (a TwitterApi) GetTrendsByPlace(v url.Values) (trends []Trend, err error) {
	trendResponse := TrendResponse{}
	response_ch := make(chan response)
	a.queryQueue <- query{BaseUrl + "/trends/place.json", v, &[]interface{}{&trendResponse}, _GET, response_ch}
	return trendResponse.Trends, (<-response_ch).err
}
