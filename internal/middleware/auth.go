package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
)

const (
	UserIDHeader = "X-User-Id"

	AccessCookieName = "access_token"
	userIDClaim      = "user_id"
)

var (
	errEmptyToken    = errors.New("empty token")
	errUserIDMissing = errors.New("user_id claim missing")
)

func Auth(jwtSecret string) fiber.Handler {
	secret := []byte(jwtSecret)

	return func(c *fiber.Ctx) error {
		token := c.Cookies(AccessCookieName)
		userID, err := parseUserID(token, secret)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    int(codes.Unauthenticated),
				"message": "unauthenticated",
				"details": []any{},
			})
		}

		c.Request().Header.Set(UserIDHeader, userID.String())
		return c.Next()
	}
}

func parseUserID(token string, secret []byte) (uuid.UUID, error) {
	if token == "" {
		return uuid.Nil, errEmptyToken
	}

	parsed, err := jwt.Parse(token, func(*jwt.Token) (any, error) {
		return secret, nil
	},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return uuid.Nil, err
	}

	userIDStr, ok := parsed.Claims.(jwt.MapClaims)[userIDClaim].(string)
	if !ok {
		return uuid.Nil, errUserIDMissing
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}
