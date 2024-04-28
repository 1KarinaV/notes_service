package database

import (
	"context"
	"github.com/google/uuid"
	"notes_service/internal/models"
)

//go:generate mockery --with-expecter --disable-version-string --case underscore --name UserRp
type UserRp interface {
	Store(ctx context.Context, user models.User) error
	Get(ctx context.Context, login string) (models.User, error)
}

//go:generate mockery --with-expecter --disable-version-string --case underscore --name NoteRp
type NoteRp interface {
	Get(ctx context.Context, userID uuid.UUID, header string) (models.Note, error)
	Store(ctx context.Context, note models.Note) error
	Update(ctx context.Context, note models.Note) error
	Delete(ctx context.Context, userID uuid.UUID, header string) error
	ListWithOptions(ctx context.Context, filterOpts models.FilterOptions) ([]models.ListedNote, error)
}
