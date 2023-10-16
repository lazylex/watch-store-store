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
	crudCreate          = "c"
	crudRead            = "r"
	crudUpdate          = "u"
	crudDelete          = "d"
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

// CheckJWT проверяет JWT токен в запросе. В случае, если токен не валидный, функция прекращает дальнейшую обработку
// запроса сервисом. Ошибка логгируется, отправителю возвращается ответ с кодом http.StatusUnauthorized. При валидном
// токене производится проверка на существование в теле токена разрешения на CRUD операцию, соответствующую http-методу
// запроса. Разрешения "c", "r", "u", "d" должны иметь значение true, чтобы считаться выданными. При несоответствии или
// отсутствии разрешения, прекращается обработка запроса, ошибка выводится в лог, а отправителю возвращается ошибка
// http.StatusForbidden.
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
			} else if errors.Is(err, jwt.ErrTokenExpired) {
				log.Error("token expired")
			} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
				log.Error("token not valid yet")
			} else {
				log.Error("couldn't handle this token:", err)
			}

			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			switch r.Method {
			case http.MethodPost:
				if claims[crudCreate] != true {
					rw.WriteHeader(http.StatusForbidden)
					log.Warn(fmt.Sprintf("trying to crudCreate without permissions. Client: %s", r.RemoteAddr))
					return
				}
			case http.MethodGet:
				if claims[crudRead] != true {
					rw.WriteHeader(http.StatusForbidden)
					log.Warn(fmt.Sprintf("trying to read without permissions. Client: %s", r.RemoteAddr))
					return
				}
			case http.MethodPut:
				if claims[crudUpdate] != true {
					rw.WriteHeader(http.StatusForbidden)
					log.Warn(fmt.Sprintf("trying to update without permissions. Client: %s", r.RemoteAddr))
					return
				}
			case http.MethodDelete:
				if claims[crudDelete] != true {
					rw.WriteHeader(http.StatusForbidden)
					log.Warn(fmt.Sprintf("trying to delete without permissions. Client: %s", r.RemoteAddr))
					return
				}
			}

		} else {
			log.Error(err.Error())
		}

		next.ServeHTTP(rw, r)
	})
}
