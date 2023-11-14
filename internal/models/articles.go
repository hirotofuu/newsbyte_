package models

import (
	"time"
)

// User describes the data for the User type.
type Article struct {
	ID         int       `json:"id"`
	Content    string    `json:"content"`
	Title      string    `json:"title"`
	Medium     int       `json:"medium"`
	UserID     int       `json:"user_id"`
	CommentOK  bool      `json:"comment_ok"`
	IsOpenFlag bool      `json:"is_open_flag"`
	TagsIn     []string  `json:"tags_in,omitempty""`
	TagsOut    string    `json:"tagss_out,omitempty""`
	Name       string    `json:"name,omitempty"`
	IdName     string    `json:"id_name,omitempty"`
	Avatar     string    `json:"avatar,omitempty"`
	GoodsCount int       `json:"good_count,omitempty"`
	IsGoodFlag int       `json:"is_good_flag,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
