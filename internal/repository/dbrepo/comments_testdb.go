package dbrepo

import (
	"errors"

	"github.com/hirotofuu/newsbyte/internal/models"
)

type TestCDBRepo struct{}

// コメント
func (m *TestCDBRepo) InsertComment(comment models.Comment) (int, error) {

	return 1, nil
}

// 記事ごとのコメント
func (m *TestCDBRepo) ArticleComments(articleID int) ([]*models.Comment, error) {

	if articleID == 1 {
		var comments []*models.Comment

		return comments, nil
	}
	return nil, errors.New("not found article")
}

// ユーザーごとのコメント
func (m *TestCDBRepo) UserComments(userID int) ([]*models.Comment, error) {

	if userID == 1 {
		var comments []*models.Comment

		return comments, nil
	}
	return nil, errors.New("not found article")
}

// コメント削除
func (m *TestCDBRepo) DeleteComment(id int) error {

	if id == 1 {
		return nil
	}
	return errors.New("not found article")
}

// コメント一つ
func (m *TestCDBRepo) OneComment(id, mainID int) (*models.Comment, error) {
	return nil, nil
}

// コメントいいね
func (m *TestCDBRepo) InsertGoodComment(id, mainID int) error {
	if id == 1 {
		return nil
	}
	return errors.New("not found article")
}

// いいね削除
func (m *TestCDBRepo) DeleteGoodComment(commentID, mainID int) error {
	if commentID == 1 {
		return nil
	}
	return errors.New("not found article")
}
