package httputil

import (
	"net/http"

	"serica-go/internal/pkg/exception"
)

type Action func(w http.ResponseWriter, r *http.Request) (any, error)

func Wrap(action Action) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := action(w, r)
		if err != nil {
			switch e := err.(type) {
			case *exception.ExError:
				e.Exception.Write(w, e.Message)
			case *exception.Exception:
				e.Write(w)
			default:
				exception.InternalServerError.Write(w, err.Error())
			}
			return
		}
		if data != nil {
			RespondJSON(w, http.StatusOK, data)
		}
	}
}
