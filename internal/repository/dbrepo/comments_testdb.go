package dbrepo

import (
	"errors"

	"github.com/hirotofuu/newsbyte/internal/models"
)

type TestCDBRepo struct{}

func (m *TestCDBRepo) InsertComment(comment models.Comment) (int, error) {

	return 1, nil
}

func (m *TestCDBRepo) ArticleComments(articleID, mainID int) ([]*models.Comment, error) {

	if articleID == 1 {
		var comments []*models.Comment

		return comments, nil
	}
	return nil, errors.New("not found article")
}

func (m *TestCDBRepo) UserComments(userID, mainID int) ([]*models.Comment, error) {

	if userID == 1 {
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

func (m *TestCDBRepo) OneComment(id, mainID int) (*models.Comment, error) {
	return nil, nil
}

func (m *TestCDBRepo) InsertGoodComment(id, mainID int) error {
	return nil
}

func (m *TestCDBRepo) DeleteGoodComment(commentID, mainID int) error {
	return nil
}
