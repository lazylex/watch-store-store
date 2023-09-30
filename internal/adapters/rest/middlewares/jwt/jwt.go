package jwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lazylex/watch-store/store/internal/logger"
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
		log := logger.AddPlaceAndRequestId(m.logger, "middlewares.jwt.CheckJWT", r)

		if len(r.Header.Get(header)) > len(requestHeaderPrefix) {
			notParsedToken = r.Header.Get(header)[len(requestHeaderPrefix):]
		} else {
			log.Error("no JWT token find")
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

		if !token.Valid {
			if errors.Is(err, jwt.ErrTokenMalformed) {
				log.Error("not a JWT token in request")
			} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
				log.Error("invalid JWT token signature")
			} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
				log.Error("not correct time in JWT token")
			} else {
				log.Error("couldn't handle this token:", err)
			}

			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(rw, r)
	})
}
