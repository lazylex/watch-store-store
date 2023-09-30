package jwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
	"net/http"
)

const (
	requestHeaderPrefix = "Bearer "
	header              = "Authorization"
)

var (
	validMethods = []string{"HS256"}
)

type MiddlewareJWT struct {
	secret []byte
	logger *slog.Logger
}

// New конструктор прослойки для проверки JSON Web Token
func New(logger *slog.Logger, secret []byte) *MiddlewareJWT {
	return &MiddlewareJWT{secret: secret, logger: logger}
}

// CheckJWT проверяет JWT токен в запросе и в случае, если токен не валидный, прекращает дальнейшую обработку запроса
// сервисом. Ошибка выводится в лог, а отправителю возвращается ответ с кодом 401
func (m *MiddlewareJWT) CheckJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		var notParsedToken string
		if len(r.Header.Get(header)) > len(requestHeaderPrefix) {
			notParsedToken = r.Header.Get(header)[len(requestHeaderPrefix):]
		} else {
			m.logger.Error("no JWT token find")
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}
		token, err := jwt.Parse(
			notParsedToken,
			func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return m.secret, nil
			},
			jwt.WithValidMethods(validMethods),
		)

		if token.Valid {
			m.logger.Info("request with valid token")
		} else {
			if errors.Is(err, jwt.ErrTokenMalformed) {
				m.logger.Error("not a JWT token in request")
			} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
				m.logger.Error("invalid JWT token signature")
			} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
				m.logger.Error("not correct time in JWT token")
			} else {
				m.logger.Error("couldn't handle this token:", err)
			}

			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(rw, r)
	})
}
