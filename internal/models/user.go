package models

type User struct {
	ID           int    `json:"id" db:"id"`
	Email        string `json:"email" db:"email"`
	Who          string `json:"who" db:"who"`
	TgChatID     int64  `json:"tg_chat_id" db:"tg_chat_id"`
	PassHash     string `json:"pass_hash" db:"pass_hash"`
	IsRegistered bool   `json:"is_registered" db:"is_registered"`
}
