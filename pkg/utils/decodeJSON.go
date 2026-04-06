package utils

import (
	"encoding/json"
	"errors"
	"net/http"
)

const (
	maxJSONBodyBytes = int64(1 << 20) // 1 MiB - обычно достаточно для auth/profile
)

// DecodeJSON декодит JSON тело запроса в dst.
//
// Меры безопасности:
// http.MaxBytesReader ограничивает размер тела, защищает от больших payload (DoS).
// DisallowUnknownFields - запрещает неизвестные поля (снижает вероятность "тихих" багов на клиентах).
// После Decode проверяем, что в body не было второго JSON объекта ("{}{}").
//
// Возвращает error - его обрабатывают handlers и возвращают 400.
func DecodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	// Ограничиваем размер тела. Если будет больше - net/http сам вернет ошибку чтения.
	r.Body = http.MaxBytesReader(w, r.Body, maxJSONBodyBytes)

	dec := json.NewDecoder(r.Body)
	// Полезно для контрактов: если клиент прислал лишнее поле - лучше упасть 400,
	// чем молча проигнорировать.
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return err
	}

	// Запрет "двух JSON объектов" в теле запроса:
	// например: {"a":1}{"b":2}
	// Это может быть как ошибка клиента, так и попытка хитрого трюка.
	if dec.More() {
		return errors.New("invalid JSON: multiple objects")
	}

	return nil
}
