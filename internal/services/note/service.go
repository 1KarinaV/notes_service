package note

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"notes_service/internal/database/pg"
	"notes_service/internal/models"
	"notes_service/internal/services"
	"time"
)

func (s *service) CreateNote(ctx context.Context, login, header, content string) error {
	userID, err := s.userSrv.GetUserID(ctx, login)
	if err != nil {
		return fmt.Errorf("cannot get user id : %w", err)
	}

	noteID, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("cannot generate uuid for new note : %w", err)
	}

	if err := s.noteRp.Store(ctx, models.Note{
		ID:      noteID,
		UserID:  userID,
		Header:  header,
		Content: content,
	}); err != nil {
		return fmt.Errorf("cannot store new note : %w", err)
	}

	return nil
}

func (s *service) GetNote(ctx context.Context, login, header string) (models.Note, error) {
	userID, err := s.userSrv.GetUserID(ctx, login)
	if err != nil {
		return models.Note{}, fmt.Errorf("cannot get user id : %w", err)
	}

	note, err := s.noteRp.Get(ctx, userID, header)
	if err != nil {
		if errors.Is(err, pg.ErrNoteNotFound) {
			return models.Note{}, services.ErrNotFound
		}

		return models.Note{}, fmt.Errorf("cannot get note : %w", err)
	}

	return note, nil
}

func (s *service) DeleteNote(ctx context.Context, login, header string) error {
	userID, err := s.userSrv.GetUserID(ctx, login)
	if err != nil {
		return fmt.Errorf("cannot get user id : %w", err)
	}

	if err := s.noteRp.Delete(ctx, userID, header); err != nil {
		return fmt.Errorf("cannot delete note : %w", err)
	}

	return nil
}

func (s *service) UpdateNote(ctx context.Context, login, header, content string) error {
	userID, err := s.userSrv.GetUserID(ctx, login)
	if err != nil {
		return fmt.Errorf("cannot get user id : %w", err)
	}

	noteDb, err := s.noteRp.Get(ctx, userID, header)
	if err != nil {
		if errors.Is(err, pg.ErrNoteNotFound) {
			return services.ErrNotFound
		}

		return fmt.Errorf("cannot get note : %w", err)
	}

	if time.Now().Sub(noteDb.CreatedAt) > MessageUpdateTimeout {
		return services.ErrNoteUpdateExceeded
	}

	if err := s.noteRp.Update(ctx, models.Note{
		UserID:  userID,
		Header:  header,
		Content: content}); err != nil {
		return fmt.Errorf("cannot update note content : %w", err)
	}

	return nil
}

func (s *service) ListWithOptions(ctx context.Context, filterOpts models.FilterOptions) ([]models.ListedNote, error) {
	return s.noteRp.ListWithOptions(ctx, filterOpts)
}
