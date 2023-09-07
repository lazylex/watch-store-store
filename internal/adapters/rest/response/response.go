package response

import (
	"errors"
	"github.com/lazylex/watch-store/store/internal/helpers/constantes/prefixes"
	"github.com/lazylex/watch-store/store/internal/ports/repository"
	"log/slog"
	"net/http"
	"strings"
)

// WriteHeaderAndLogAboutErr записывает заголовок ответа сервера, соответствующий переданной ошибке. К примеру, при
// отсутствующей записи, записывает заголовок ответа http.StatusNotFound. Также текст ошибки записывается в лог. При
// отсутсвии ошибки, функция ничего не выполняет
func WriteHeaderAndLogAboutErr(w http.ResponseWriter, logger *slog.Logger, err error) {
	if err == nil {
		return
	}

	switch {
	case strings.HasPrefix(err.Error(), prefixes.RequestErrorsPrefix):
		w.WriteHeader(http.StatusBadRequest)
	case strings.HasPrefix(err.Error(), prefixes.DTOErrorsPrefix):
		w.WriteHeader(http.StatusBadRequest)
	case errors.Is(err, repository.ErrNoRecord):
		w.WriteHeader(http.StatusNotFound)
	case errors.Is(err, repository.ErrTimeout):
		w.WriteHeader(http.StatusRequestTimeout)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	logger.Warn(err.Error())
}

// WriteHeaderAndLogAboutBadRequest записывает заголовок ответа http.StatusBadRequest и переданную ошибку в лог
func WriteHeaderAndLogAboutBadRequest(w http.ResponseWriter, logger *slog.Logger, err error) {
	w.WriteHeader(http.StatusBadRequest)
	logger.Warn(err.Error())
}
