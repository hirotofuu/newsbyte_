package dbrepo

import (
	"errors"

	"github.com/hirotofuu/newsbyte/internal/models"
)

type TestCDBRepo struct{}

func (m *TestDBRepo) InsertComment(comment models.Comment) (int, error) {

	return 1, nil
}

func (m *TestDBRepo) ArticleComments(articleID int) ([]*models.Comment, error) {
	var comments []*models.Comment

	return comments, nil
}

func (m *TestDBRepo) UserComments(userID int) ([]*models.Comment, error) {
	var comments []*models.Comment

	return comments, nil
}

func (m *TestDBRepo) DeleteComment(id int) error {

	if id == 1 {
		return nil
	}
	return errors.New("not found article")
}
