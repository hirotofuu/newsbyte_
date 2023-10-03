package repository

import (
	"github.com/hirotofuu/newsbyte/internal/models"
)

type ArticleDatabaserepo interface {
	InsertArticle(article models.Article) (int, error)
	AllArticles() ([]*models.Article, error)
	UserArticles(userID int) ([]*models.Article, error)
	WorkArticles(work string) ([]*models.Article, error)
	OneArticle(id int) (*models.Article, error)
	DeleteArticle(id int) error
}
