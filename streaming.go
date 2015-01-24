package anaconda

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/dustin/go-jsonpointer"
)

const (
	BaseUrlUserStream = "https://userstream.twitter.com/1.1"
	BaseUrlSiteStream = "https://sitestream.twitter.com/1.1"
	BaseUrlStream     = "https://stream.twitter.com/1.1"
)

// messages

type StatusDeletionNotice struct {
	Id        int64  `json:"id"`
	IdStr     string `json:"id_str"`
	UserId    int64  `json:"user_id"`
	UserIdStr string `json:"user_id_str"`
}
type statusDeletionNotice struct {
	Delete *struct {
		Status *StatusDeletionNotice `json:"status"`
	} `json:"delete"`
}

type LocationDeletionNotice struct {
	UserId          int64  `json:"user_id"`
	UserIdStr       string `json:"user_id_str"`
	UpToStatusId    int64  `json:"up_to_status_id"`
	UpToStatusIdStr string `json:"up_to_status_id_str"`
}
type locationDeletionNotice struct {
	ScrubGeo *LocationDeletionNotice `json:"scrub_geo"`
}

type LimitNotice struct {
	Track int64 `json:"track"`
}
type limitNotice struct {
	Limit *LimitNotice `json:"limit"`
}

type StatusWithheldNotice struct {
	Id                  int64    `json:"id"`
	UserId              int64    `json:"user_id"`
	WithheldInCountries []string `json:"withheld_in_countries"`
}
type statusWithheldNotice struct {
	StatusWithheld *StatusWithheldNotice `json:"status_withheld"`
}

type UserWithheldNotice struct {
	Id                  int64    `json:"id"`
	WithheldInCountries []string `json:"withheld_in_countries"`
}
type userWithheldNotice struct {
	UserWithheld *UserWithheldNotice `json:"user_withheld"`
}

type DisconnectMessage struct {
	Code       int64  `json:"code"`
	StreamName string `json:"stream_name"`
	Reason     string `json:"reason"`
}
type disconnectMessage struct {
	Disconnect *DisconnectMessage `json:"disconnect"`
}

type StallWarning struct {
	Code        string `json:"code"`
	Message     string `json:"message"`
	PercentFull int64  `json:"percent_full"`
}
type stallWarning struct {
	Warning *StallWarning `json:"warning"`
}

type FriendsList []int64
type friendsList struct {
	Friends *FriendsList `json:"friends"`
}

type streamDirectMessage struct {
	DirectMessage *DirectMessage `json:"direct_message"`
}

type Event struct {
	Target    *User  `json:"target"`
	Source    *User  `json:"source"`
	Event     string `json:"event"`
	CreatedAt string `json:"created_at"`
}

// TODO: List struct is not defined
// type EventList struct {
// 	TargetObject *List `json:"target_object"`
// }

type EventTweet struct {
	Event
	TargetObject *Tweet `json:"target_object"`
}

type TooManyFollow struct {
	Warning *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		UserId  int64  `json:"user_id"`
	} `json:"warning"`
}

// TODO: Site Stream messages. I cant test.

type Stream struct {
	response *http.Response
	C        chan interface{}
}

func (s Stream) Close() error {
	close(s.C)
	return s.response.Body.Close()
}

func (s Stream) listen() {
	go func() {
		defer s.Close()

		scanner := bufio.NewScanner(s.response.Body)
		for {
			if ok := scanner.Scan(); !ok {
				break
			}
			// TODO: DRY
			j := scanner.Bytes()
			if scanner.Text() == "" {
				continue
			} else if o := new(Tweet); jsonAsStruct(j, "/source", o) {
				s.C <- *o
			} else if o := new(statusDeletionNotice); jsonAsStruct(j, "/delete", o) {
				s.C <- *o.Delete.Status
			} else if o := new(locationDeletionNotice); jsonAsStruct(j, "/scrub_geo", o) {
				s.C <- *o.ScrubGeo
			} else if o := new(limitNotice); jsonAsStruct(j, "/limit", o) {
				s.C <- *o.Limit
			} else if o := new(statusWithheldNotice); jsonAsStruct(j, "/status_withheld", o) {
				s.C <- *o.StatusWithheld
			} else if o := new(userWithheldNotice); jsonAsStruct(j, "/user_withheld", o) {
				s.C <- *o.UserWithheld
			} else if o := new(disconnectMessage); jsonAsStruct(j, "/disconnect", o) {
				s.C <- *o.Disconnect
			} else if o := new(stallWarning); jsonAsStruct(j, "/warning", o) {
				s.C <- *o.Warning
			} else if o := new(friendsList); jsonAsStruct(j, "/friends", o) {
				s.C <- *o.Friends
			} else if o := new(streamDirectMessage); jsonAsStruct(j, "/direct_message", o) {
				s.C <- *o.DirectMessage
			} else if o := new(EventTweet); jsonAsStruct(j, "/target_object/source", o) {
				s.C <- *o
			} else if o := new(Event); jsonAsStruct(j, "/target_object", o) {
				s.C <- *o
			}
		}
	}()
}

func (a TwitterApi) newStream(urlStr string, v url.Values, method int) (stream Stream, err error) {
	var resp *http.Response
	switch method {
	case _GET:
		resp, err = oauthClient.Get(a.HttpClient, a.Credentials, urlStr, v)
	case _POST:
		resp, err = oauthClient.Post(a.HttpClient, a.Credentials, urlStr, v)
	default:
		return stream, fmt.Errorf("HTTP method not yet supported")
	}
	if err != nil {
		return
	}

	stream = Stream{resp, make(chan interface{})}
	go stream.listen()
	return
}

func (a TwitterApi) UserStream(v url.Values) (stream Stream, err error) {
	return a.newStream(BaseUrlUserStream+"/user.json", v, _GET)
}

func (a TwitterApi) PublicStreamSample(v url.Values) (stream Stream, err error) {
	return a.newStream(BaseUrlStream+"/statuses/sample.json", v, _GET)
}

// XXX: To use this API authority is requied. but I dont have this. I cant test.
func (a TwitterApi) PublicStreamFirehose(v url.Values) (stream Stream, err error) {
	return a.newStream(BaseUrlStream+"/statuses/firehose.json", v, _GET)
}

// XXX: PublicStream(Track|Follow|Locations) func is needed?
func (a TwitterApi) PublicStreamFilter(v url.Values) (stream Stream, err error) {
	return a.newStream(BaseUrlStream+"/statuses/filter.json", v, _POST)
}

// XXX: To use this API authority is requied. but I dont have this. I cant test.
func (a TwitterApi) SiteStream(v url.Values) (stream Stream, err error) {
	return a.newStream(BaseUrlSiteStream+"/site.json", v, _GET)
}

func jsonAsStruct(j []byte, path string, obj interface{}) (res bool) {
	if v, _ := jsonpointer.Find(j, path); v == nil {
		return false
	}
	err := json.Unmarshal(j, obj)
	return err == nil
}
