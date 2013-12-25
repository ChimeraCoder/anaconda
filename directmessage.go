package anaconda

type DirectMessage struct {
	CreatedAt           string          `json:"created_at"`
	Entities            TwitterEntities `json:"entities"`
	Id                  int64           `json:"id"`
	IdStr               string          `json:"id_str"`
	Recipient           TwitterUser     `json:"recipient"`
	RecipientId         int64           `json:"recipient_id"`
	RecipientScreenName string          `json:"recipient_screen_name"`
	Sender              TwitterUser     `json:"sender"`
	SenderId            int64           `json:"sender_id"`
	SenderScreenName    string          `json:"sender_screen_name"`
	Text                string          `json:"text"`
}
