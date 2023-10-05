package models

import (
	"time"
)

// User describes the data for the User type.
type Article struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	Work      string    `json:"work"`
	MainImg   string    `json:"main_img,omitempty"`
	Medium    int       `json:"medium"`
	UserID    int       `json:"user_id"`
	CommentOK bool      `json:"comment_ok"`
	Name      string    `json:"name,omitempty"`
	Avatar    string    `json:"avatar,omitempty"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
