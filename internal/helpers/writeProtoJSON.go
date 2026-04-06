package helpers

import (
	"gateWay/internal/DTO"
	"google.golang.org/protobuf/proto"
	"net/http"
)

// WriteProtoJSON сериализует protobuf message в JSON и пишет HTTP ответ
// Если marshal сломался - это почти всегда внутренняя ошибка (500)
func WriteProtoJSON(w http.ResponseWriter, statusCode int, msg proto.Message) {

	b, err := DTO.ProtoJSON.Marshal(msg)
	if err != nil {
		// fallback: простой текстовый 500
		http.Error(w, "failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_, _ = w.Write(b)
}
