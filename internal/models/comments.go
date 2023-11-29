package models

import "time"

type Comment struct {
	ID            int       `json:"id"`
	Comment       string    `json:"comment"`
	UserID        int       `json:"user_id"`
	ArticleID     int       `json:"article_id"`
	Name          string    `json:"name,omitempty"`
	Avatar        string    `json:"avatar,omitempty"`
	GoodsCount    int       `json:"good_count,omitempty"`
	ArticleTitle  string    `json:"article_title,omitempty"`
	ArticleUserId string    `json:"article_user_id,omitempty"`
	IsGoodFlag    int       `json:"is_good_flag,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
