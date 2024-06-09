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
	urlPermissions string        // Адрес сервиса безопасности для получения нумерованных разрешений
	urlLogin       string        // Адрес сервиса безопасности для входа в учетную запись и получения токена
	username       string        // Логин данного приложения в сервисе безопасности
	password       string        // Пароль данного приложения в сервисе безопасности
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

	urlPermissions := fmt.Sprintf("%s://%s/get-numbered-permissions?service=store", cfg.Protocol, cfg.Server)
	urlLogin := fmt.Sprintf("%s://%s/login", cfg.Protocol, cfg.Server)

	return &Secure{attempts: attempts,
		urlPermissions: urlPermissions,
		urlLogin:       urlLogin,
		requestTimeout: cfg.RequestTimeout,
		username:       cfg.Username,
		password:       cfg.Password}
}

// login получает токен сессии в микросервисе secure.
func (s *Secure) login() (string, error) {
	log := slog.Default().With(logger.OPLabel, "secure.login")

	type resultJSON struct {
		Token string `json:"token"`
	}

	var response *http.Response
	var result resultJSON

	ctx, cancel := context.WithTimeout(context.Background(), s.requestTimeout)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, s.urlLogin, nil)
	if err != nil {
		return "", err
	}

	request.SetBasicAuth(s.username, s.password)

	client := http.DefaultClient
	response, err = client.Do(request)

	if err != nil {
		return "", err
	}

	if response.StatusCode == http.StatusUnauthorized {
		return "", ErrUnauthorized
	}

	var bytes []byte
	if response.StatusCode == http.StatusOK {
		bytes, err = responseBodyBytes(response, 36)
		if err != nil {
			return "", err
		}

		err = json.Unmarshal(bytes, &result)
		if err != nil {
			return "", err
		}
		log.Info("successfully login to secure service")
	}

	return result.Token, err
}

// MustGetPermissionsNumbers получает список нумерованных разрешений и отправляет его в канал nameNumbersChan. Если все
// попытки (количество указывается в конфигурации) оказались неудачными, приложение завершает работу. Перед очередной
// попыткой выдерживается пауза, которая каждый раз увеличивается на одну секунду.
func (s *Secure) MustGetPermissionsNumbers(nameNumbersChan chan dto.NameNumber) {
	var result []dto.NameNumber
	var err error
	defer close(nameNumbersChan)

	log := slog.Default().With(logger.OPLabel, "secure.MustGetPermissionsNumbers")

	if s.tokens.secure == "" {
		s.tokens.secure, err = s.login()
	}

	if err != nil {
		log.Error(fmt.Errorf("failed to obtain permissions (reason: %w)", err).Error())
		os.Exit(1)
	}

	for attempt := 0; attempt < s.attempts; attempt++ {
		if errors.Is(err, ErrUnauthorized) {
			s.tokens.secure, err = s.login()
		}

		result, err = s.getPermissionsNumbers(s.tokens.secure, s.urlPermissions)

		if err == nil {
			for _, nameNumber := range result {
				nameNumbersChan <- nameNumber
			}
			return
		}
		log.Warn(fmt.Sprintf("failed to obtain permissions (attempt %d)", attempt+1))
		time.Sleep(time.Duration(attempt) * time.Second)
	}

	if err != nil {
		log.Error(fmt.Errorf("failed to obtain permissions (reason: %w)", err).Error())
		os.Exit(1)
	}
}

// getPermissionsNumbers получение списка нумерованных разрешений.
func (s *Secure) getPermissionsNumbers(token, url string) ([]dto.NameNumber, error) {
	var result []dto.NameNumber
	var response *http.Response

	log := slog.Default().With(logger.OPLabel, "secure.getPermissionsNumbers")

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

	var bytes []byte
	if response.StatusCode == http.StatusOK {
		bytes, err = responseBodyBytes(response, 1024)

		err = json.Unmarshal(bytes, &result)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

// responseBodyBytes возвращает слайс байт из тела ответа.
func responseBodyBytes(response *http.Response, allocateBytes int) ([]byte, error) {
	var bodyBytes []byte
	var err error
	var n int

	defer func() {
		_ = response.Body.Close()
	}()

	bytes := make([]byte, allocateBytes)

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

	return bodyBytes, nil
}
