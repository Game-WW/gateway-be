package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

const (
	UserAgentHeader    = "X-User-Agent"
	RefreshTokenHeader = "X-Refresh-Token"

	RefreshCookieName = "refresh_token"

	forwardedForHeader = "X-Forwarded-For"
)

// Headers готовит заголовки, которые proxy переложит в gRPC metadata.
// Всё, что клиент мог прислать в них сам, перезаписывается или удаляется.
func Headers() fiber.Handler {
	return func(c *fiber.Ctx) error {
		h := &c.Request().Header

		// X-User-Id ставит только Auth — на публичных ручках его быть не должно.
		h.Del(UserIDHeader)

		// grpc-gateway дописывает реальный IP к X-Forwarded-For, а users-be берёт первый адрес —
		// без удаления им был бы адрес, который прислал клиент.
		h.Del(forwardedForHeader)

		// User-Agent в gRPC перезаписывает сам клиент, поэтому передаём копию.
		setOrDel(h, UserAgentHeader, c.Get(fiber.HeaderUserAgent))
		setOrDel(h, RefreshTokenHeader, c.Cookies(RefreshCookieName))

		return c.Next()
	}
}

func setOrDel(h *fasthttp.RequestHeader, key, value string) {
	if value == "" {
		h.Del(key)
		return
	}
	h.Set(key, value)
}
