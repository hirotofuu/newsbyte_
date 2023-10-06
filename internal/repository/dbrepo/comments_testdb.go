package dbrepo

import (
	"errors"

	"github.com/hirotofuu/newsbyte/internal/models"
)

type TestCDBRepo struct{}

func (m *TestCDBRepo) InsertComment(comment models.Comment) (int, error) {

	return 1, nil
}

func (m *TestCDBRepo) ArticleComments(articleID int) ([]*models.Comment, error) {

	if (articleID == 1) {
		var comments []*models.Comment
	
		return comments, nil
	}
	return nil, errors.New("not found article")
}

func (m *TestCDBRepo) UserComments(userID int) ([]*models.Comment, error) {

	if (userID == 1) {
		var comments []*models.Comment
	
		return comments, nil
	}
	return nil, errors.New("not found article")
}

func (m *TestCDBRepo) DeleteComment(id int) error {

	if id == 1 {
		return nil
	}
	return errors.New("not found article")
}
