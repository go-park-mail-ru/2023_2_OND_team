package v1

import (
	"net/http"
)

func (h *HandlerHTTP) CreateNewBoard(w http.ResponseWriter, r *http.Request) {
	/*
		1. достаем user_id из контекста запроса
		2. парсим тело в json, проверяем парсинг
		3. вызываем UseCase, проверяем ошибку
		4. формируем ответ: заголовок content-type, body ответа
	*/
}
