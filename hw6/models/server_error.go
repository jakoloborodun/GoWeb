package models

// ErrorModel - Ошибка отвечаемая сервером
type ErrorModel struct {
	Code     int
	Err      string
	Desc     string
	Internal interface{}
}
