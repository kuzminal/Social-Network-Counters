package models

type UserInfo struct {
	Id         string `json:"id" msgpack:"Id"`
	FirstName  string `json:"first_name" msgpack:"FirstName"`
	SecondName string `json:"second_name" msgpack:"SecondName"`
	Age        int    `json:"age" msgpack:"Age"`
	Birthdate  string `json:"birthdate" msgpack:"Birthdate"`
	Biography  string `json:"biography" msgpack:"Biography"`
	City       string `json:"city" msgpack:"City"`
	Password   string `json:"-" msgpack:"Password"`
}

type AuthInfo struct {
	Id       string `json:"id"`
	Password string `json:"password"`
}

type UserSession struct {
	Id        string `json:"id" msgpack:"Id"`
	UserId    string `json:"userId" msgpack:"user_id"`
	Token     string `json:"token" msgpack:"token"`
	CreatedAt uint64 `json:"createdAt" msgpack:"created_at"`
}
