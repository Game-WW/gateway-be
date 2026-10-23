package proxy

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Game-WW/gateway-be/internal/middleware"
	"github.com/golang-jwt/jwt/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/proto"
)

const (
	accessTokenMetaKey  = "x-access-token"
	refreshTokenMetaKey = "x-refresh-token"
	logoutMetaKey       = "x-logout"

	refreshTokenTTL = 30 * 24 * time.Hour
)

// forwardCookies превращает токены из metadata ответа сервиса в cookie.
func forwardCookies(ctx context.Context, w http.ResponseWriter, _ proto.Message) error {
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		return nil
	}

	if len(md.HeaderMD.Get(logoutMetaKey)) > 0 {
		clearCookie(w, middleware.AccessCookieName)
		clearCookie(w, middleware.RefreshCookieName)
		return nil
	}

	if v := md.HeaderMD.Get(accessTokenMetaKey); len(v) > 0 {
		ttl, err := accessTokenTTL(v[0])
		if err != nil {
			return err
		}
		setCookie(w, middleware.AccessCookieName, v[0], ttl)
	}
	if v := md.HeaderMD.Get(refreshTokenMetaKey); len(v) > 0 {
		setCookie(w, middleware.RefreshCookieName, v[0], refreshTokenTTL)
	}

	return nil
}

// accessTokenTTL берёт срок жизни cookie из exp самого токена.
func accessTokenTTL(token string) (time.Duration, error) {
	parsed, _, err := jwt.NewParser().ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		return 0, fmt.Errorf("proxy: parse access token: %w", err)
	}

	exp, err := parsed.Claims.GetExpirationTime()
	if err != nil || exp == nil {
		return 0, fmt.Errorf("proxy: access token without exp")
	}

	return time.Until(exp.Time), nil
}

func setCookie(w http.ResponseWriter, name, value string, ttl time.Duration) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(ttl.Seconds()),
	})
}

func clearCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}
