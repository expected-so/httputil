package httputil

import (
	"encoding/json"
	"net/http"
)

type (
	Response interface {
		Write(w http.ResponseWriter) error
	}

	responseWriter func(w http.ResponseWriter) error
)

var (
	JSONInternalServerError = NewJSONResponse(
		http.StatusInternalServerError,
		map[string]string{"message": "Internal server error."},
	)

	JSONUnauthorizedError = NewJSONResponse(
		http.StatusUnauthorized,
		map[string]string{"message": "Unauthorized."},
	)

	JSONBadRequestError = NewJSONResponse(
		http.StatusBadRequest,
		map[string]string{"message": "Bad request."},
	)
)

func (r responseWriter) Write(w http.ResponseWriter) error {
	return r(w)
}

func NewJSONResponse[T any](status int, body T) Response {
	return responseWriter(func(w http.ResponseWriter) error {
		w.WriteHeader(status)
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		_, err = w.Write(b)
		return err
	})
}

func NewEmptyResponse(status int) Response {
	return responseWriter(func(w http.ResponseWriter) error {
		w.WriteHeader(status)
		return nil
	})
}

func NewRedirectResponse(status int, url string) Response {
	return responseWriter(func(w http.ResponseWriter) error {
		w.Header().Set("Location", url)
		w.WriteHeader(status)
		return nil
	})
}
