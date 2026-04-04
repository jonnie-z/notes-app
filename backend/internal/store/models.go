package store

type Note struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

type NoteRepository interface {
	GetAll() ([]Note, error)
	Create(text string) (Note, error)
	Update(id int, text string) (Note, error)
	Delete(id int) error
}