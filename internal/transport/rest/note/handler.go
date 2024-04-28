package note

import (
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"net/http"
	"notes_service/internal/models"
	"notes_service/internal/services"
	"notes_service/pkg/web"
	"strconv"
	"time"
)

func (h *Handler) CreateNote(w http.ResponseWriter, r *http.Request) {
	login, err := h.authService.GetClaims(r)
	if err != nil {
		render.Render(w, r, web.ErrRender(err, http.StatusUnauthorized))
		return
	}

	var note Note
	if err := render.DecodeJSON(r.Body, &note); err != nil {
		render.Render(w, r, web.ErrRender(err, http.StatusBadRequest))
		return
	}

	if err := h.noteService.CreateNote(r.Context(), login, note.Header, note.Content); err != nil {
		render.Render(w, r, web.ErrRender(err, http.StatusInternalServerError))
		return
	}

	createdNote, err := h.noteService.GetNote(r.Context(), login, note.Header)
	if err != nil {
		render.Render(w, r, web.ErrRender(err, http.StatusNotFound))
		return
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, Note{
		Login:     login,
		Header:    createdNote.Header,
		Content:   createdNote.Content,
		CreatedAt: createdNote.CreatedAt,
	})
}

func (h *Handler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	login, err := h.authService.GetClaims(r)
	if err != nil {
		render.Render(w, r, web.ErrRender(err, http.StatusUnauthorized))
		return
	}

	var note Note
	if err := render.DecodeJSON(r.Body, &note); err != nil {
		render.Render(w, r, web.ErrRender(err, http.StatusBadRequest))
		return
	}

	if err := h.noteService.DeleteNote(r.Context(), login, note.Header); err != nil {
		render.Render(w, r, web.ErrRender(err, http.StatusInternalServerError))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	login, err := h.authService.GetClaims(r)
	if err != nil {
		render.Render(w, r, web.ErrRender(err, http.StatusUnauthorized))
		return
	}

	var note Note
	if err := render.DecodeJSON(r.Body, &note); err != nil {
		render.Render(w, r, web.ErrRender(err, http.StatusBadRequest))
		return
	}

	if err := h.noteService.UpdateNote(r.Context(), login, note.Header, note.Content); err != nil {
		var status int
		switch {
		case errors.Is(err, services.ErrNotFound):
			status = http.StatusNotFound
		case errors.Is(err, services.ErrNoteUpdateExceeded):
			status = http.StatusBadRequest
		default:
			status = http.StatusInternalServerError
		}

		render.Render(w, r, web.ErrRender(err, status))

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) ListNotes(w http.ResponseWriter, r *http.Request) {
	filterOpts, err := h.parseFilterOpts(r)
	if err != nil {
		render.Render(w, r, web.ErrRender(err, http.StatusBadRequest))
		return
	}

	noteList, err := h.noteService.ListWithOptions(r.Context(), filterOpts)

	if err != nil {
		render.Render(w, r, web.ErrRender(err, http.StatusInternalServerError))
		return
	}

	login, err := h.authService.GetClaims(r)
	if err != nil && !errors.Is(err, services.ErrEmptyClaims) {
		render.Render(w, r, web.ErrRender(err, http.StatusUnauthorized))
		return
	}

	resultNotesList := make([]ListedNote, 0, len(noteList))

	for _, note := range noteList {
		listedNote := ListedNote{
			Login:   note.Login,
			Header:  note.Header,
			Content: note.Content,
		}

		if login == note.Login {
			listedNote.MyNote = true
		}

		resultNotesList = append(resultNotesList, listedNote)
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, resultNotesList)
}

func (h *Handler) parseFilterOpts(r *http.Request) (models.FilterOptions, error) {
	filterOpts := models.FilterOptions{
		Page:  PageDefault,
		Limit: LimitDefault,
	}

	page := r.URL.Query().Get(Page)
	if page != "" {
		pageNum, err := strconv.Atoi(page)
		if err != nil {
			return models.FilterOptions{}, fmt.Errorf("cannot parse page : %w", err)
		}

		if pageNum < 0 {
			return models.FilterOptions{}, ErrInvalidPageParam
		}

		filterOpts.Page = uint64(pageNum)
	}

	limit := r.URL.Query().Get(Limit)
	if limit != "" {
		limitNum, err := strconv.Atoi(limit)
		if err != nil {
			return models.FilterOptions{}, fmt.Errorf("cannot parse limit : %w", err)
		}

		if limitNum < 0 {
			return models.FilterOptions{}, ErrInvalidLimitParam
		}

		filterOpts.Limit = uint64(limitNum)
	}

	author := r.URL.Query().Get(Author)
	if author != "" {
		filterOpts.Author = &author
	}

	startTimeStr := r.URL.Query().Get(StartTime)
	if startTimeStr != "" {
		startTime, err := time.Parse(Layout, startTimeStr)
		if err != nil {
			return models.FilterOptions{}, fmt.Errorf("cannot parse start time : %w", err)
		}

		filterOpts.StartDate = &startTime
	}

	endTimeStr := r.URL.Query().Get(EndTime)
	if endTimeStr != "" {
		endTime, err := time.Parse(Layout, endTimeStr)
		if err != nil {
			return models.FilterOptions{}, fmt.Errorf("cannot parse start time : %w", err)
		}

		filterOpts.EndDate = &endTime
	}

	if filterOpts.StartDate != nil && filterOpts.EndDate != nil {
		if filterOpts.StartDate.After(*filterOpts.EndDate) {
			return filterOpts, fmt.Errorf("invalid time range")
		}
	}

	return filterOpts, nil
}
