package rest

import "net/http"

/*
JSONMsg можно возвращать из хендлеров, чтобы вернуть JSON-ответ с данными или ошибкой.
плюсы:
	- ответ и ошибка в одном объекте
	- можно задать HTTP-код ответа
	- заготовленные константы для ответов OK и Accepted
*/

var (
	OK       = NewJSONMsg("OK").WithHTTPCode(http.StatusOK)
	Accepted = NewJSONMsg("Accepted").WithHTTPCode(http.StatusAccepted)
)

type JSONMsg[T any] struct {
	Data      T      `json:"data,omitempty"`
	ErrorDesc string `json:"errorDesc,omitempty"`
	ErrorCode int    `json:"errorCode"`
	httpCode  int
}

func NewJSONMsg[T any](data T) *JSONMsg[T] {
	return &JSONMsg[T]{Data: data}
}

type Event struct {
	Name string
}

func (it JSONMsg[T]) WithHTTPCode(code int) JSONMsg[T] {
	it.httpCode = code
	return it
}

func (it JSONMsg[T]) HTTPCode() int {
	if it.httpCode != 0 {
		return it.httpCode
	}

	return http.StatusOK
}
