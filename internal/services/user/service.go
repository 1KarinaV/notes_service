package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"notes_service/internal/database/pg"
	"notes_service/internal/models"
	"notes_service/internal/services"
)

func (s *service) CreateNewUser(ctx context.Context, login, password string) error {
	passwordHash, err := s.hashService.Hash(password)
	if err != nil {
		return fmt.Errorf("cannot create hash : %w", err)
	}

	userUUID, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("cannot generate id for user : %w", err)
	}

	if err := s.userRp.Store(ctx, models.User{
		ID:       userUUID,
		Login:    login,
		Password: passwordHash,
	}); err != nil {
		return fmt.Errorf("cannot store user with login %s : %w", login, err)
	}

	return nil
}

func (s *service) CompareUserPassword(ctx context.Context, user models.User) (bool, error) {
	usr, err := s.userRp.Get(ctx, user.Login)
	if err != nil {
		if errors.Is(err, pg.ErrUserNotFound) {
			return false, services.ErrNotFound
		}
		return false, fmt.Errorf("error while get user : %w", err)
	}

	return s.hashService.CompareHash(user.Password, usr.Password), nil
}

func (s *service) GetUserID(ctx context.Context, login string) (uuid.UUID, error) {
	user, err := s.userRp.Get(ctx, login)
	if err != nil {
		if errors.Is(err, pg.ErrUserNotFound) {
			return uuid.Nil, services.ErrNotFound
		}

		return uuid.Nil, fmt.Errorf("error while get user : %w", err)
	}

	return user.ID, nil
}
