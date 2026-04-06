package auth

import (
	"fmt"
	DTO "gateWay/internal/DTO/Auth"
	"gateWay/internal/helpers"
	"gateWay/pkg/utils"
	authv1 "github.com/EvgGo/proto/proto/gen/go/sso"
	"google.golang.org/grpc/status"
	"log/slog"
	"net/http"
	"strings"
)

// RegisterHandler обрабатывает POST /auth/register:
//
// 1) Декодит JSON тело
// 2) Делает базовую валидацию (email/password не пустые)
// 3) Нормализует email (lowercase + trim)
// 4) Готовит ctx с timeout и metadata (authorization/request-id)
// 5) Вызывает gRPC Register
// 6) Возвращает protobuf-ответ в JSON
func RegisterHandler(log *slog.Logger, c authv1.AuthClient) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		reqLog := log.With(
			"op", "RegisterHandler",
			"method", r.Method,
			"path", r.URL.Path,
			"remote", r.RemoteAddr,
			"request_id", r.Header.Get("X-Request-ID"),
		)

		var in DTO.RegisterJSON
		if err := utils.DecodeJSON(w, r, &in); err != nil {
			reqLog.Warn("invalid JSON body",
				"err", err,
				"content_type", r.Header.Get("Content-Type"),
				"content_length", r.ContentLength,
			)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid JSON body")
			return
		}

		email := strings.ToLower(strings.TrimSpace(in.Email))
		password := strings.TrimSpace(in.Password)

		if email == "" || password == "" {

			reqLog.Warn("validation failed",
				"email_empty", email == "",
				"password_empty", password == "",
			)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "email and password are required")
			return
		}

		firstName := strings.TrimSpace(in.FirstName)
		lastName := strings.TrimSpace(in.LastName)
		phone := strings.TrimSpace(in.Phone)

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()
		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		resp, err := c.Register(ctx, &authv1.RegisterRequest{
			Email:     email,
			Password:  in.Password, // не логируем
			FirstName: firstName,
			LastName:  lastName,
			Phone:     phone,
		})
		if err != nil {
			st, _ := status.FromError(err)
			reqLog.Warn("Register failed",
				"email", maskEmail(email),
				"grpc_code", st.Code().String(),
				"grpc_message", st.Message(),
				"err", err,
			)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		reqLog.Info("Register succeeded",
			"email", maskEmail(email),
			"has_first_name", firstName != "",
			"has_last_name", lastName != "",
			"has_phone", phone != "",
		)

		log.Debug(fmt.Sprintf("На регистрацию получили %v", resp))

		helpers.WriteJSON(w, http.StatusCreated, map[string]string{
			"status": "ok",
		})
	}
}

func maskEmail(email string) string {
	email = strings.TrimSpace(email)
	at := strings.IndexByte(email, '@')
	if at <= 1 {
		if email == "" {
			return ""
		}
		return "***"
	}
	return email[:1] + "***" + email[at:]
}
