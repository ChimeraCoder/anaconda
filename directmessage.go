package anaconda

type DirectMessage struct {
	CreatedAt           string   `json:"created_at"`
	Entities            Entities `json:"entities"`
	Id                  int64    `json:"id"`
	IdStr               string   `json:"id_str"`
	Recipient           User     `json:"recipient"`
	RecipientId         int64    `json:"recipient_id"`
	RecipientScreenName string   `json:"recipient_screen_name"`
	Sender              User     `json:"sender"`
	SenderId            int64    `json:"sender_id"`
	SenderScreenName    string   `json:"sender_screen_name"`
	Text                string   `json:"text"`
}

type DMEventData struct {
	DMEvent *DMEvent `json:"event"`
}

type DMEventList struct {
	NextCursor string    `json:"next_cursor"`
	DMEvents   []DMEvent `json:"events"`
}

type DMEvent struct {
	Type             string `json:"type"`
	Id               string `json:"id"`
	CreatedTimestamp string `json:"created_timestamp"`

	InitiatedVia *struct {
		TweetId          string `json:"tweet_id"`
		WelcomeMessageId string `json:"welcome_message_id"`
	} `json:"initiated_via"`

	MessageCreate *MessageCreate `json:"message_create"`
}

type MessageCreate struct {
	Target struct {
		RecipientId string `json:"recipient_id"`
	} `json:"target"`
	SenderId    string       `json:"sender_id"`
	SourceAppId string       `json:"source_app_id"`
	MessageData *MessageData `json:"message_data"`
}

type MessageData struct {
	Text               string   `json:"text"`
	Entities           Entities `json:"entities"`
	QuickReplyResponse *struct {
		Type     string `json:"type"`
		Metadata string `json:"metadata"`
	} `json:"quick_reply_response"`
	Attachment *struct {
		Type  string      `json:"type"`
		Media EntityMedia `json:"media"`
	} `json:"attachment"`
}
