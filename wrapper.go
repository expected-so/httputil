package httputil

import (
	"net/http"
)

type Wrapper struct {
	SupportedErrors map[error]func(r *http.Request) Response
	HandleError     func(err any) Response
}

func (wr Wrapper) handleError(r *http.Request, rawErr any) Response {
	if err, ok := rawErr.(error); ok && wr.SupportedErrors != nil {
		if handler, ok := wr.SupportedErrors[err]; ok {
			return handler(r)
		}
	}
	if wr.HandleError != nil {
		return wr.HandleError(rawErr)
	}
	return nil
}

func (wr Wrapper) Wrap(handler func(r *http.Request) (Response, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		var res Response

		defer func() {
			if recoverErr := recover(); err != nil {
				res = wr.handleError(r, recoverErr)
			} else if err != nil {
				res = wr.handleError(r, err)
			}

			if res != nil {
				_ = res.Write(w)
			}
		}()

		res, err = handler(r)
	}
}
