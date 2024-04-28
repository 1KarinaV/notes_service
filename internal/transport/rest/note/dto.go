package note

import "time"

type Note struct {
	Login     string    `json:"login,omitempty"`
	Header    string    `validate:"min=10,max=255" json:"header,omitempty"`
	Content   string    `validate:"max=10500" json:"content,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type ListedNote struct {
	Login   string `json:"login"`
	Header  string `json:"header"`
	Content string `json:"content"`
	MyNote  bool   `json:"my_note,omitempty"`
}
