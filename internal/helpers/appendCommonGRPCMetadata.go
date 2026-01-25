package helpers

import (
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc/metadata"
	"net/http"
)

// AppendCommonGRPCMetadata добавляет общие metadata, которые полезны для каждого gRPC вызова:
//
//	"authorization" - прокидываем Bearer токен как есть (если он есть)
//	"x-request-id"  - request id из chi middleware для корреляции логов/трейсов
//
//	 мы используем NewOutgoingContext - создаем новый контекст с metadata,
//	 не мутируя исходный
func AppendCommonGRPCMetadata(ctx context.Context, r *http.Request) context.Context {

	md := metadata.New(nil)

	// Прокинем access token (если есть)
	if auth := BearerFromRequest(r); auth != "" {
		md.Append("authorization", auth)
	}

	// Прокинем request id (из chi middleware)
	if rid := middleware.GetReqID(r.Context()); rid != "" {
		md.Append("x-request-id", rid)
	}

	return metadata.NewOutgoingContext(ctx, md)
}
