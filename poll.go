package anaconda

import (
	"fmt"
	"net/url"
)

//{
//  "request": {
//    "params": {
//      "first_choice": "East",
//      "name": "best coast poll",
//      "second_choice": "West",
//      "media_key": "13_950589518557540353",
//      "duration_in_minutes": 10080
//    }
//  },
//  "data": {
//    "video_poster_height": "9",
//    "name": "best coast poll",
//    "start_time": "2018-01-09T04:51:34Z",
//    "first_choice": "East",
//    "video_height": "9",
//    "video_url": "https://video.twimg.com/amplify_video/vmap/950589518557540353.vmap",
//    "content_duration_seconds": "8",
//    "second_choice": "West",
//    "end_time": "2018-01-16T04:51:34Z",
//    "id": "57i77",
//    "video_width": "16",
//    "video_hls_url": "https://video.twimg.com/amplify_video/950589518557540353/vid/1280x720/BRkAhPxFoBREIaFA.mp4",
//    "created_at": "2018-01-09T04:51:34Z",
//    "duration_in_minutes": "10080",
//    "card_uri": "card://950590850777497601",
//    "updated_at": "2018-01-09T04:51:34Z",
//    "video_poster_url": "https://pbs.twimg.com/amplify_video_thumb/950589518557540353/img/nZ1vX_MXYqmvbsXP.jpg",
//    "video_poster_width": "16",
//    "deleted": false,
//    "card_type": "VIDEO_POLLS"
//  }
//}
type PollCard struct {
	Request PollRequest `json:"request"`
	Data    PollData    `json:"data"`
}

type PollData struct {
	Name              string `json:"name"`
	FirstChoice       string `json:"first_choice"`
	SecondChoice      string `json:"second_choice"`
	ThirdChoice       string `json:"third_choice"`
	FourthChoice      string `json:"fourthChoice"`
	MediaKey          string `json:"media_key"`
	DurationInMinutes int    `json:"duration_in_minutes"`
	CardType          string `json:"card_type"`
	CardURI           string `json:"card_uri"`
}

type PollParams struct {
	FirstChoice       string `json:"first_choice"`
	SecondChoice      string `json:"second_choice"`
	ThirdChoice       string `json:"third_choice"`
	FourthChoice      string `json:"fourthChoice"`
	MediaKey          string `json:"media_key"`
	DurationInMinutes int    `json:"duration_in_minutes"`
}

type PollRequest struct {
	Params PollParams `json:"params"`
}

//CreatePollCard will create a new poll card that can be attached to a tweet
func (a TwitterApi) CreatePollCard(choices []string, accountID string, durationInMinutes int, v url.Values) (card PollCard, err error) {
	// TODO @kris-nova we can support media here which would be cool because it's in beta and you can't do that from the web UI
	var firstChoice, secondChoice, thirdChoice, fourthChoice string
	choiceCount := len(choices)
	if choiceCount == 1 {
		firstChoice = choices[0]
		if len(firstChoice) > 25 {
			return card, fmt.Errorf("choice longer than 25 characters")
		}
	}
	if choiceCount == 2 {
		secondChoice = choices[1]
		if len(secondChoice) > 25 {
			return card, fmt.Errorf("choice longer than 25 characters")
		}
	}
	if choiceCount == 3 {
		thirdChoice = choices[2]
		if len(thirdChoice) > 25 {
			return card, fmt.Errorf("choice longer than 25 characters")
		}
	}
	if choiceCount == 4 {
		fourthChoice = choices[3]
		if len(fourthChoice) > 25 {
			return card, fmt.Errorf("choice longer than 25 characters")
		}
	}
	v = cleanValues(v)
	v.Set("first_choice", firstChoice)
	v.Set("second_choice", secondChoice)
	v.Set("third_choice", thirdChoice)
	v.Set("fourth_choice", fourthChoice)
	v.Set("account_id", accountID)
	v.Set("duration_in_minutes", fmt.Sprintf("%d", durationInMinutes))
	response_ch := make(chan response)
	a.queryQueue <- query{AddBaseUrl + fmt.Sprintf("/account/%s/cards/poll.json", accountID), v, &card, _POST, response_ch}
	return card, (<-response_ch).err
}
