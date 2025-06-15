package models

import "encoding/gob"

// Регистрирует тип []SessionResume для gob-сериализации,
// чтобы его можно было сохранять в сессии (например, через gin-contrib/sessions).
func init() {
	gob.Register([]SessionResume{})
}

type SessionResume struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}
