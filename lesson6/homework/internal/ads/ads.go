package ads

import "fmt"

type Ad struct {
	ID        int64
	Title     string
	Text      string
	AuthorID  int64
	Published bool
}

func (a *Ad) String() string {
	return fmt.Sprintf(
		"<Ad id=%d authorID=%d published=%v title=`%s` text=`%s`>",
		a.ID,
		a.AuthorID,
		a.Published,
		a.Title,
		a.Text,
	)
}

func New(id int64, title string, text string, authorID int64) *Ad {
	return &Ad{
		ID:       id,
		Title:    title,
		Text:     text,
		AuthorID: authorID,
	}
}
