package helpers

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

// WriteGRPCError конвертит gRPC error в HTTP ошибку
//
//	берем code и превращаем в HTTP status, а message отдаем клиенту
//
// - Если error не gRPC status (неожиданный), считаем это 500
func WriteGRPCError(w http.ResponseWriter, r *http.Request, err error) {

	st, ok := status.FromError(err)
	if !ok {
		WriteAPIError(w, r, http.StatusInternalServerError, "internal error")
		return
	}

	httpCode := grpcCodeToHTTP(st.Code())

	// Клиенту - короткое сообщение; детали остаются в логах gateway и/или внутри gRPC сервисов
	WriteAPIError(w, r, httpCode, st.Message())
}

// grpcCodeToHTTP мапит gRPC codes на HTTP коды
func grpcCodeToHTTP(c codes.Code) int {
	switch c {
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.FailedPrecondition:
		return http.StatusPreconditionFailed
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}
