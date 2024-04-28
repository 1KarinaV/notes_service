package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"gopkg.in/validator.v2"
	"net/http"
	"notes_service/internal/models"
	"notes_service/internal/services"
	"notes_service/pkg/web"
)

func (h *Handler) AuthCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var credentials Credentials
		if err := render.DecodeJSON(r.Body, &credentials); err != nil {
			render.Render(w, r, web.ErrRender(err, http.StatusBadRequest))
			return
		}

		if err := validator.Validate(credentials); err != nil {
			render.Render(w, r, web.ErrRender(err, http.StatusBadRequest))
			return
		}

		ctx := context.WithValue(r.Context(), Creds, credentials)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	credentials, ok := r.Context().Value(Creds).(Credentials)
	if !ok {
		render.Render(w, r, web.ErrRender(ErrInvalidCast, http.StatusInternalServerError))
		return
	}

	rs, err := h.userService.CompareUserPassword(r.Context(), models.User{
		Login:    credentials.Login,
		Password: credentials.Password,
	})
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			render.Render(w, r, web.ErrRender(fmt.Errorf("cannot find user with login %s", credentials.Login), http.StatusNotFound))
		} else {
			render.Render(w, r, web.ErrRender(err, http.StatusBadRequest))
		}
		return
	}

	if !rs {
		render.Render(w, r, web.ErrRender(ErrInvalidPassword, http.StatusBadRequest))
		return
	}

	token, err := h.authService.MakeToken(credentials.Login)
	if err != nil {
		render.Render(w, r, web.ErrRender(err, http.StatusInternalServerError))
		return
	}

	w.Header().Add(AuthHeader, fmt.Sprintf("Bearer %s", token))
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(AuthHeader, "")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	credentials := r.Context().Value(Creds).(Credentials)
	if err := h.userService.CreateNewUser(r.Context(), credentials.Login, credentials.Password); err != nil {
		render.Render(w, r, web.ErrRender(err, http.StatusBadRequest))
		return
	}

	w.WriteHeader(http.StatusOK)
}
