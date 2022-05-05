package anaconda

import (
	"fmt"
	"net/url"
)

func (a TwitterApi) GetCollectionListByUserId(userId int64, v url.Values) (result CollectionListResult, err error) {
	v = cleanValues(v)
	v.Set("user_id", fmt.Sprintf("%d", userId))
	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/collections/list.json", v, &result, _GET, response_ch}
	return result, (<-response_ch).err
}

func (a TwitterApi) GetCollectionListByScreenName(screenName string, v url.Values) (result CollectionListResult, err error) {
	v = cleanValues(v)
	v.Set("screen_name", screenName)
	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/collections/list.json", v, &result, _GET, response_ch}
	return result, (<-response_ch).err
}

func (a TwitterApi) GetCollectionShow(id string, v url.Values) (result CollectionShowResult, err error) {
	v = cleanValues(v)
	v.Set("id", id)
	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/collections/show.json", v, &result, _GET, response_ch}
	return result, (<-response_ch).err
}

func (a TwitterApi) GetCollectionEntries(id string, v url.Values) (result CollectionEntriesResult, err error) {
	v = cleanValues(v)
	v.Set("id", id)
	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/collections/entries.json", v, &result, _GET, response_ch}
	return result, (<-response_ch).err
}

func (a TwitterApi) CreateCollection(name string, v url.Values) (result CollectionShowResult, err error) {
	v = cleanValues(v)
	v.Set("name", name)
	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/collections/create.json", v, &result, _POST, response_ch}
	return result, (<-response_ch).err
}

func (a TwitterApi) UpdateCollection(id string, v url.Values) (result CollectionShowResult, err error) {
	v = cleanValues(v)
	v.Set("id", id)
	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/collections/update.json", v, &result, _POST, response_ch}
	return result, (<-response_ch).err
}

func (a TwitterApi) DestroyCollection(id string, v url.Values) (result CollectionDestroyResult, err error) {
	v = cleanValues(v)
	v.Set("id", id)
	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/collections/destroy.json", v, &result, _POST, response_ch}
	return result, (<-response_ch).err
}

func (a TwitterApi) AddEntryToCollection(id string, tweetId int64, v url.Values) (result CollectionEntryAddResult, err error) {
	v = cleanValues(v)
	v.Set("id", id)
	v.Set("tweet_id", fmt.Sprintf("%d", tweetId))
	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/collections/entries/add.json", v, &result, _POST, response_ch}
	return result, (<-response_ch).err
}

func (a TwitterApi) RemoveEntryFromCollection(id string, tweetId int64, v url.Values) (result CollectionEntryAddResult, err error) {
	v = cleanValues(v)
	v.Set("id", id)
	v.Set("tweet_id", fmt.Sprintf("%d", tweetId))
	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/collections/entries/remove.json", v, &result, _POST, response_ch}
	return result, (<-response_ch).err
}

func (a TwitterApi) MoveEntryFromCollection(id string, tweetId, relativeTo int64, v url.Values) (result CollectionEntryAddResult, err error) {
	v = cleanValues(v)
	v.Set("id", id)
	v.Set("tweet_id", fmt.Sprintf("%d", tweetId))
	v.Set("relative_to", fmt.Sprintf("%d", relativeTo))
	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/collections/entries/move.json", v, &result, _POST, response_ch}
	return result, (<-response_ch).err
}
