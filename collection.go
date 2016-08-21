package anaconda

type Collection struct {
	Name           string
	UserId         string `json:"user_id"`
	CollectionUrl  string `json:"collection_url"`
	Description    string
	Url            string
	Visibility     string
	TimelineOrder  string `json:"timeline_order"`
	CollectionType string `json:"collection_type"`
}

type CollectionListResult struct {
	Objects struct {
		Users     map[string]User
		Timelines map[string]Collection
	}
	Response struct {
		Results []struct {
			TimelineId string `json:"timeline_id"`
		}
		Cursors struct {
			NextCursor string `json:"next_cursor"`
		}
	}
}

// GET collections/show
// Also used by POST collections/create
// Also used by POST collections/update
type CollectionShowResult struct {
	Objects struct {
		Users     map[string]User
		Timelines map[string]Collection
	}
	Response struct {
		TimelineId string `json:"timeline_id"`
	}
}

type CollectionEntriesResult struct {
	Objects struct {
		Timelines map[string]Collection
	}
	Tweets   []Tweet
	Response struct {
		Position struct {
			MaxPosition  string `json:"max_position"`
			MinPosition  string `json:"min_position"`
			WasTruncated bool   `json:"was_truncated"`
		}
		Timeline []struct {
			FeatureContext string `json:"feature_context"`
			Tweet          struct {
				Id        string
				SortIndex string `json:"sort_index"`
			}
		}
		TimelineId string `json:"timeline_id"`
	}
}

type CollectionDestroyResult struct {
	Destroyed bool
}

// POST collections/entries/add
// Also used by POST collections/entries/remove
// Also used by POST collections/entries/move
type CollectionEntryAddResult struct {
	Objects  struct{}
	Response struct {
		Errors []struct {
			Change struct {
				Op      string
				TweetId string `json:"tweet_id"`
			}
			Reason string
		}
	}
}
