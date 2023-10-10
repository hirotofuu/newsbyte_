package repository

import (
	"github.com/hirotofuu/newsbyte/internal/models"
)

type Databaserepo interface {
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
	InsertUser(user models.User) (int, error)
	AllUsers() ([]*models.User, error)
	InsertFollow(id, mainID int) error
	DeleteFollow(id, mainID int) error
}

type ArticleDatabaserepo interface {
	InsertArticle(article models.Article) (int, error)
	AllArticles() ([]*models.Article, error)
	UserArticles(userID int) ([]*models.Article, error)
	WorkArticles(work string) ([]*models.Article, error)
	OneArticle(id, mainID int) (*models.Article, error)
	DeleteArticle(id int) error
	InsertGoodArticle(id, mainID int) error
	DeleteGoodArticle(articleID, mainID int) error
}

type CommentDatabaserepo interface {
	InsertComment(comment models.Comment) (int, error)
	ArticleComments(articleID, mainID int) ([]*models.Comment, error)
	UserComments(userID, mainID int) ([]*models.Comment, error)
	OneComment(id, mainID int) (*models.Comment, error)
	DeleteComment(id int) error
	InsertGoodComment(id, mainID int) error
	DeleteGoodComment(commentID, mainID int) error
}
