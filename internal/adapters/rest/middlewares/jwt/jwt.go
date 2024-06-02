package jwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lazylex/watch-store/store/internal/logger"
	"log/slog"
	"net/http"
	"reflect"
)

const (
	requestHeaderPrefix = "Bearer "
	header              = "Authorization"
)

var (
	validMethods = []string{"HS256"}
)

type MiddlewareJWT struct {
	secret      []byte         // Секретный ключ, которым должны быть подписаны валидные token-ы
	permissions map[string]int // Карта номеров разрешений для путей, где ключ - путь, а его значение - номер разрешения
}

// New конструктор прослойки для проверки JSON Web Token
func New(secret []byte, permissions map[string]int) *MiddlewareJWT {
	return &MiddlewareJWT{secret: secret, permissions: permissions}
}

// CheckJWT проверяет JWT токен в запросе. В случае, если токен не валидный, функция прекращает дальнейшую обработку
// запроса сервисом. Ошибка заносится в лог, отправителю возвращается ответ с кодом http.StatusUnauthorized.
func (m *MiddlewareJWT) CheckJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		var notParsedToken string
		log := logger.AddPlaceAndRequestId(slog.Default(), "adapters.rest.middlewares.jwt.CheckJWT", r)

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
				log.Warn("not a JWT token in request")
			} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
				log.Warn("invalid JWT token signature")
			} else if errors.Is(err, jwt.ErrTokenExpired) {
				log.Warn("token expired")
			} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
				log.Warn("token not valid yet")
			} else {
				log.Warn("couldn't handle this token:", err)
			}

			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		if err = m.checkPermissions(token, r.URL.Path); err != nil {
			rw.WriteHeader(http.StatusForbidden)
			log.Warn(err.Error())
			return
		}

		next.ServeHTTP(rw, r)
	})
}

// checkPermissions проверяет наличие в token-е номера разрешения, соответствующего переданному url. При нахождении
// такого номера возвращается nil, в любом другом случае - возвращается ошибка. Разрешения в token-е должны содержаться
// по ключу perm в полезной нагрузке (claims).
func (m *MiddlewareJWT) checkPermissions(token *jwt.Token, url string) error {
	if _, ok := m.permissions[url]; !ok {
		return fmt.Errorf("no such url: %v", url)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		var (
			permissions []int
			rawPerm     interface{}
			ok          bool
		)

		if rawPerm, ok = claims["perm"]; !ok {
			return fmt.Errorf("no permissions claims in token")
		}

		if reflect.TypeOf(rawPerm).Kind() == reflect.Slice {
			s := reflect.ValueOf(rawPerm)
			for i := 0; i < s.Len(); i++ {
				element := s.Index(i)
				if val, noProblem := element.Interface().(float64); noProblem {
					permissions = append(permissions, int(val))
				}
			}
		} else {
			return fmt.Errorf("can't read permissions claims")
		}

		if len(permissions) == 0 {
			return fmt.Errorf("empty permissions list")
		}

		for _, v := range permissions {
			if v == m.permissions[url] {
				return nil
			}
		}
	}

	return fmt.Errorf("no permissions for this path")
}
