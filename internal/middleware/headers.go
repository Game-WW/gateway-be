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
func Headers() fiber.Handler {
	return func(c *fiber.Ctx) error {
		h := &c.Request().Header

		h.Del(UserIDHeader)
		h.Del(forwardedForHeader)

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
