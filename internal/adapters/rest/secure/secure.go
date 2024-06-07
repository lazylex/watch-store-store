package secure

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lazylex/watch-store/store/internal/config"
	"github.com/lazylex/watch-store/store/internal/dto"
	"github.com/lazylex/watch-store/store/internal/logger"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
)

type Secure struct {
	tokens         Tokens        // Токены
	attempts       int           // Количество попыток обращения к стороннему сервису
	url            string        // Адрес сервиса безопасности
	requestTimeout time.Duration // Таймаут запроса
}

// Tokens предназначен для хранения токенов, необходимых для работы приложения.
type Tokens struct {
	secure string // Токен для обращения к микросервису secure
}

// New возвращает указатель на структуру, предназначенную для работы с токенами и разрешениями.
func New(cfg *config.Secure) *Secure {
	attempts := cfg.Attempts
	if attempts == 0 {
		attempts = 3
	}

	protocol := cfg.Protocol
	if protocol == "" {
		protocol = "http"
	}

	url := fmt.Sprintf("%s://%s/get-numbered-permissions?service=store", cfg.Protocol, cfg.Server)

	return &Secure{attempts: attempts, url: url, requestTimeout: cfg.RequestTimeout}
}

// login получает токен сессии в микросервисе secure.
func login() (string, error) {
	// TODO implement
	slog.Debug("login not implemented")
	return "", nil
}

// MustGetPermissionsNumbers получение списка нумерованных разрешений. Если все попытки (количество указывается в
// конфигурации) оказались неудачными, приложение завершает работу. Перед очередной попыткой выдерживается пауза,
// которая каждый раз увеличивается на одну секунду.
func (s *Secure) MustGetPermissionsNumbers() ([]dto.NameNumber, error) {
	var result []dto.NameNumber
	var err error

	log := slog.Default().With(logger.OPLabel, "secure.MustGetPermissionsNumbers")

	if s.tokens.secure == "" {
		s.tokens.secure, err = login()
	}

	if err != nil {
		return nil, err
	}

	for attempt := 0; attempt < s.attempts; attempt++ {
		if errors.Is(err, ErrUnauthorized) {
			s.tokens.secure, err = login()
		}

		result, err = s.getPermissionsNumbers(s.tokens.secure, s.url)

		if err == nil {
			return result, nil
		}

		log.Warn(fmt.Sprintf("failed to obtain permissions (attempt %d)", attempt+1))
		time.Sleep(time.Duration(attempt) * time.Second)
	}

	if err != nil {
		log.Error(fmt.Errorf("failed to obtain permissions (reason: %w)", err).Error())
		os.Exit(1)
	}

	return result, err
}

// getPermissionsNumbers получение списка нумерованных разрешений.
func (s *Secure) getPermissionsNumbers(token, url string) ([]dto.NameNumber, error) {
	var result []dto.NameNumber
	var response *http.Response

	log := slog.Default().With("op", "secure.getPermissionsNumbers")

	ctx, cancel := context.WithTimeout(context.Background(), s.requestTimeout)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	client := http.DefaultClient
	response, err = client.Do(request)

	if err != nil {
		return nil, err
	}

	if response.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	}

	if response.StatusCode == http.StatusOK {
		var bodyBytes []byte
		var n int

		defer func() {
			_ = response.Body.Close()
		}()

		bytes := make([]byte, 1024)

		for {
			bytes = bytes[:cap(bytes)]
			n, err = response.Body.Read(bytes)

			if err != nil {
				if err == io.EOF {
					bodyBytes = append(bodyBytes, bytes[:n]...)
					break
				}
				return nil, err
			}
			bodyBytes = append(bodyBytes, bytes[:n]...)
		}

		err = json.Unmarshal(bodyBytes, &result)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
