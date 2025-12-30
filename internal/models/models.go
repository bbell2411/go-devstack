package models

import "time"

type Users struct {
	ID         int
	Name       string
	Pin        int
	Created_at time.Time
}

type Snippets struct {
	ID         int
	User_id    int
	Snippet    string
	Language   string
	Deleted_at *time.Time
	Created_at time.Time
	Updated_at *time.Time
}
