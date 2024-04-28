package web

import (
	"github.com/go-chi/render"
	"net/http"
)

type ErrResponse struct {
	ErrMessage     string `json:"err_message"`
	HTTPStatusCode int    `json:"status"`
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrRender(err error, statusCode int) render.Renderer {
	return &ErrResponse{
		ErrMessage:     err.Error(),
		HTTPStatusCode: statusCode,
	}
}
