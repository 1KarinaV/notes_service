package models

import (
	"github.com/google/uuid"
	"time"
)

type Note struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	Header    string    `db:"header"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type ListedNote struct {
	Login   string `db:"login"`
	Header  string `db:"header"`
	Content string `db:"content"`
}
