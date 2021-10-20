package main

type Chat struct {
	ChatId int `json:"id"`
}

type User struct {
	UserId   int    `json:"id"`
	IsBot    bool   `json:"is_bot"`
	Name     string `json:"first_name"`
	UserName string `json:"username"`
}

type Message struct {
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
	User User   `json:"from"`
}

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type RestResponse struct {
	Status bool     `json:"ok"`
	Result []Update `json:"result"`
}

type BotMessage struct {
	ChatId int    `json:"chat_id"`
	Text   string `json:"text"`
}
