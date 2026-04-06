package helpers

import (
	"context"
	"net/http"
	"time"
)

const (
	DefaultTimeout = 12 * time.Second
)

// CtxWithOutgoingMeta создает контекст для исходящего вызова (в gRPC),
// гарантируя разумный timeout, НО не ломая уже установленный дедлайн
//
// Зачем так:
//   - chi.middleware.Timeout может уже поставить дедлайн на весь HTTP запрос.
//     Если мы поверх сделаем context.WithTimeout - получится двойной дедлайн,
//     и можно неожиданно укоротить время обработки
//   - Если дедлайна нет - ставим defaultTimeout для защиты от зависаний на gRPC вызовах
//
// Возвращает (ctx, cancel):
// - cancel нужен, чтобы освободить ресурсы таймера, если дедлайн поставили здесь
// - Если дедлайн уже был - возвращаем no-op cancel, чтобы caller мог безопасно defer cancel()
func CtxWithOutgoingMeta(r *http.Request) (context.Context, context.CancelFunc) {

	ctx := r.Context()

	// Если chi.Timeout уже поставил дедлайн - не дублируем
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, DefaultTimeout)
		return ctx, cancel
	}

	// no-op cancel - удобно для единообразного defer cancel()
	return ctx, func() {}
}
