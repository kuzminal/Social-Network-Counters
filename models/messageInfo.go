package models

type Message struct {
	Id        string `json:"-" msgpack:"Id"`
	Text      string `json:"text" msgpack:"Text"`
	FromUser  string `json:"from" msgpack:"FromUser"`
	ToUser    string `json:"to" msgpack:"ToUser"`
	ChatId    string `json:"-" msgpack:"ChatId"`
	CreatedAt string `json:"-" msgpack:"CreatedAt"`
	IsRead    bool   `json:"isRead" msgpack:"IsRead"`
}
