package dbrepo

import (
	"github.com/hirotofuu/newsbyte/internal/models"
)

type TestADBRepo struct{}

func (m *TestADBRepo) InsertArticle(article models.Article) (int, error) {

	return 1, nil
}
