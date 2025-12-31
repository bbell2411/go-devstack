package models

type Users struct {
	ID        int
	Name      string
	Pin       int
	CreatedAt string
}

type Snippets struct {
	ID        int
	UserID    int
	Snippet   string
	Language  string
	DeletedAt *string
	CreatedAt string
	UpdatedAt *string
}
