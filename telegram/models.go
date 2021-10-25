package telegram

/**
| ============== Types ============== |
*/
type Chat struct {
	Id    int    `json:"id"`
	Type  string `json:"type"`
	Title string `json:"title"`
}

type User struct {
	Id        int    `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	UserName  string `json:"username"`
}

type Message struct {
	Id   int    `json:"message_id"`
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
	User User   `json:"from"`
}

type Update struct {
	UpdateId      int     `json:"update_id"`
	Message       Message `json:"message"`
	EditedMessage Message `json:"edited_message"`
}

type InputMediaPhoto struct {
	Type    string `json:"type"`  // val = photo
	Media   string `json:"media"` // a file_id that exists on the Telegram OR  an HTTP URL
	Caption string `json:"caption,omitempty"`
}

/**
| ============== Custom Types ============== |
*/
type RestResponse struct {
	Status bool     `json:"ok"`
	Result []Update `json:"result"`
}

/**
| ============== Structures to be sent ============== |
*/
type OutMessage struct {
	ChatId        int    `json:"chat_id"`
	Text          string `json:"text"`
	ReplayToMsgId int    `json:"reply_to_message_id,omitempty"`
}

type OutPhoto struct {
	ChatId          int    `json:"chat_id"`
	Photo           string `json:"photo"`
	Caption         string `json:"caption,omitempty"`
	CaptionEntities string `json:"caption_entities,omitempty"`
	ReplayToMsgId   int    `json:"reply_to_message_id,omitempty"`
}

type OutVideo struct {
	ChatId          int    `json:"chat_id"`
	Video           string `json:"video"`
	Caption         string `json:"caption,omitempty"`
	CaptionEntities string `json:"caption_entities,omitempty"`
	ReplayToMsgId   int    `json:"reply_to_message_id,omitempty"`
}

type OutMediaPhotoGroup struct {
	ChatId        int               `json:"chat_id"`
	Media         []InputMediaPhoto `json:"media"`
	ReplayToMsgId int               `json:"reply_to_message_id,omitempty"`
}
