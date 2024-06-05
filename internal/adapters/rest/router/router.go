package router

import (
	"github.com/go-chi/chi"
	"net/http"
)

var paths []string
var mux *chi.Mux

func init() {
	mux = chi.NewRouter()
}

// Mux возвращает мультиплексор.
func Mux() *chi.Mux {
	return mux
}

// AssignPathToHandler добавляет пусть к списку используемых и прикрепляет его к переданному четвертым аргументом
// обработчику.
func AssignPathToHandler(path, method string, mux *chi.Mux, handler func(http.ResponseWriter, *http.Request)) {
	switch method {
	case http.MethodGet:
		mux.Get(path, handler)
	case http.MethodPost:
		mux.Post(path, handler)
	case http.MethodPut:
		mux.Put(path, handler)
	case http.MethodPatch:
		mux.Patch(path, handler)
	case http.MethodDelete:
		mux.Delete(path, handler)
	default:
		return
	}

	paths = append(paths, path)
}

// ExistentPaths возвращает все зарегистрированные для обработчиков адреса.
func ExistentPaths() []string {
	return paths
}

// IsExistPath возвращает true, если в приложении используется передаваемый путь. Иначе - false.
func IsExistPath(path string) bool {
	for _, p := range ExistentPaths() {
		if p == path {
			return true
		}
	}

	return false
}
