package services

import (
	"context"
	"github.com/google/uuid"
	"net/http"
	"notes_service/internal/models"
)

type UserService interface {
	CreateNewUser(ctx context.Context, login, password string) error
	GetUserID(ctx context.Context, login string) (uuid.UUID, error)
	CompareUserPassword(ctx context.Context, user models.User) (bool, error)
}

type HashService interface {
	Hash(input string) (string, error)
	CompareHash(input string, hash string) bool
}

type AuthService interface {
	MakeToken(email string) (string, error)
	GetClaims(request *http.Request) (string, error)
}

type NoteService interface {
	CreateNote(ctx context.Context, login, header, content string) error
	GetNote(ctx context.Context, login, header string) (models.Note, error)
	DeleteNote(ctx context.Context, login, header string) error
	UpdateNote(ctx context.Context, login, header, content string) error
	ListWithOptions(ctx context.Context, filterOpts models.FilterOptions) ([]models.ListedNote, error)
}
