package note

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"notes_service/internal/database/pg"
	"notes_service/internal/models"
	"time"
)

func (r *repository) Store(ctx context.Context, note models.Note) error {
	query, args, err := squirrel.Insert(pg.Notes).
		Columns(pg.IdColumn,
			pg.UserIDColumn,
			pg.HeaderColumn,
			pg.ContentColumn,
			pg.CreatedAtColumn,
			pg.UpdatedAtColumn).
		Values(note.ID, note.UserID, note.Header, note.Content, time.Now(), time.Now()).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("cannot create query : %w", err)
	}

	if _, err := r.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("error while exec query : %w", err)
	}

	return nil
}

func (r *repository) Update(ctx context.Context, note models.Note) error {
	query, args, err := squirrel.Update(pg.Notes).
		Set(pg.ContentColumn, note.Content).
		Set(pg.UpdatedAtColumn, time.Now()).
		Where(squirrel.Eq{pg.UserIDColumn: note.UserID, pg.HeaderColumn: note.Header}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("cannot create query : %w", err)
	}

	if _, err := r.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("error while exec query : %w", err)
	}

	return nil
}

func (r *repository) Get(ctx context.Context, userID uuid.UUID, header string) (models.Note, error) {
	query, args, err := squirrel.Select(pg.IdColumn,
		pg.UserIDColumn,
		pg.HeaderColumn,
		pg.ContentColumn,
		pg.CreatedAtColumn,
		pg.UpdatedAtColumn).
		From(pg.Notes).
		Where(squirrel.Eq{pg.UserIDColumn: userID, pg.HeaderColumn: header}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return models.Note{}, fmt.Errorf("cannot create query : %w", err)
	}

	var note models.Note
	if err := r.GetContext(ctx, &note, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Note{}, pg.ErrNoteNotFound
		}

		return models.Note{}, fmt.Errorf("error while exec query : %w", err)
	}

	return note, nil
}

func (r *repository) Delete(ctx context.Context, userID uuid.UUID, header string) error {
	query, args, err := squirrel.Delete(pg.Notes).
		Where(squirrel.Eq{pg.UserIDColumn: userID, pg.HeaderColumn: header}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("cannot create query : %w", err)
	}

	if _, err := r.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("error while exec query : %w", err)
	}

	return nil
}

func (r *repository) ListWithOptions(ctx context.Context, filterOpts models.FilterOptions) ([]models.ListedNote, error) {
	offset := filterOpts.Page * filterOpts.Limit

	queryTmp := squirrel.Select(pg.LoginColumn, pg.HeaderColumn, pg.ContentColumn).
		From(pg.Notes).
		Join(fmt.Sprintf("%s on %s.%s = %s.id", pg.Users, pg.Notes, pg.UserIDColumn, pg.Users))

	if filterOpts.Author != nil {
		queryTmp = queryTmp.Where(squirrel.Eq{fmt.Sprintf("%s.%s", pg.Users, pg.LoginColumn): filterOpts.Author})
	}

	if filterOpts.StartDate != nil && filterOpts.EndDate != nil {
		queryTmp = queryTmp.Where(fmt.Sprintf("date(%s) BETWEEN date('%s') AND date('%s')", pg.CreatedAtColumn, filterOpts.StartDate.String()[:10], filterOpts.EndDate.String()[:10]))
	} else if filterOpts.StartDate != nil {
		queryTmp = queryTmp.Where(squirrel.Eq{fmt.Sprintf("date(%s)", pg.CreatedAtColumn): *filterOpts.StartDate})
	}

	queryTmp = queryTmp.Offset(offset).
		Limit(filterOpts.Limit).
		OrderBy(fmt.Sprintf("%s desc", pg.CreatedAtColumn))

	query, args, err := queryTmp.PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("cannot create query : %w", err)
	}

	fmt.Println(query)

	var listedNotes []models.ListedNote
	if err := r.SelectContext(ctx, &listedNotes, query, args...); err != nil {
		return nil, fmt.Errorf("error while exec query : %w", err)
	}

	return listedNotes, nil
}
