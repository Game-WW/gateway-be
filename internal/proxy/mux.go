package proxy

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	userspb "github.com/Game-WW/gateway-be/gen/users"
	"github.com/Game-WW/gateway-be/internal/middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

// Сколько ждём ответа сервиса, потом отдаём 504.
const requestTimeout = 5 * time.Second

// HTTP-заголовок запроса → ключ gRPC metadata.
var incomingHeaders = map[string]string{
	middleware.UserIDHeader:       "x-user-id",
	middleware.UserAgentHeader:    "x-user-agent",
	middleware.RefreshTokenHeader: refreshTokenMetaKey,
	"Authorization":               "authorization",
}

func NewMux(ctx context.Context, usersConn *grpc.ClientConn, cookieSecure bool) (http.Handler, error) {
	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(headerMatcher),
		runtime.WithOutgoingHeaderMatcher(dropAll),
		runtime.WithOutgoingTrailerMatcher(dropAll),
		runtime.WithForwardResponseOption(forwardCookies(cookieSecure)),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{UseProtoNames: true},
		}),
	)

	if err := userspb.RegisterUsersServiceHandler(ctx, mux, usersConn); err != nil {
		return nil, fmt.Errorf("proxy: register users: %w", err)
	}

	return withTimeout(mux), nil
}

func withTimeout(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
		defer cancel()

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// В metadata пропускаем только заголовки из incomingHeaders.
func headerMatcher(key string) (string, bool) {
	for header, metaKey := range incomingHeaders {
		if strings.EqualFold(key, header) {
			return metaKey, true
		}
	}
	return "", false
}

func dropAll(string) (string, bool) {
	return "", false
}
