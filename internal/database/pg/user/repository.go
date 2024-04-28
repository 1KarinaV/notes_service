package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"notes_service/internal/database/pg"
	"notes_service/internal/models"
)

func (r *repository) Store(ctx context.Context, user models.User) error {
	query, args, err := squirrel.Insert(pg.Users).
		Columns(pg.IdColumn, pg.LoginColumn, pg.PasswordColumn).
		Values(user.ID, user.Login, user.Password).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("cannot create query : %w", err)
	}

	if _, err = r.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("error while exec query : %w", err)
	}

	return nil
}

func (r *repository) Get(ctx context.Context, login string) (models.User, error) {
	query, args, err := squirrel.Select(pg.IdColumn, pg.LoginColumn, pg.PasswordColumn).
		From(pg.Users).
		Where(squirrel.Eq{pg.LoginColumn: login}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return models.User{}, fmt.Errorf("cannot create query : %w", err)
	}

	var user models.User
	if err := r.GetContext(ctx, &user, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, pg.ErrUserNotFound
		}

		return models.User{}, fmt.Errorf("error while exec query : %w", err)
	}

	return user, nil
}
