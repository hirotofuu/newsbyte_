package dbrepo

import (
	"errors"
	"time"

	"github.com/hirotofuu/newsbyte/internal/models"
)

type TestADBRepo struct{}

func (m *TestADBRepo) InsertArticle(article models.Article) (int, error) {

	return 1, nil
}

func (m *TestADBRepo) AllArticles() ([]*models.Article, error) {
	var articles []*models.Article

	return articles, nil
}

func (m *TestADBRepo) UserArticles(userID int) ([]*models.Article, error) {
	var articles []*models.Article

	return articles, nil
}

func (m *TestADBRepo) UserSaveArticles(userID int) ([]*models.Article, error) {
	var articles []*models.Article

	return articles, nil
}

func (m *TestADBRepo) WorkArticles(work string) ([]*models.Article, error) {
	var articles []*models.Article

	return articles, nil
}

func (m *TestADBRepo) DeleteArticle(id int) error {

	if id == 1 {
		return nil
	}
	return errors.New("not found article")
}

func (m *TestADBRepo) OneArticle(id, mainID int) (*models.Article, error) {
	var testArticle models.Article

	tag := "ワンピース"
	var tags []string
	tags = append(tags, tag)

	testArticle.ID = id
	testArticle.Title = "you know say"
	testArticle.Content = "tomorrow tomorrow i love yeah tomorrow"
	testArticle.TagsIn = tags
	testArticle.Medium = 1
	testArticle.UserID = 1
	testArticle.CommentOK = true
	testArticle.CreatedAt = time.Now()
	testArticle.UpdatedAt = time.Now()
	testArticle.Name = "hiroto"
	testArticle.Avatar = "http://hello"
	testArticle.IsGoodFlag = 1
	testArticle.GoodsCount = 0

	if id == 1 {
		return &testArticle, nil
	}

	return nil, errors.New("not found article")

}

func (m *TestADBRepo) InsertGoodArticle(id, mainID int) error {
	if id == 1 {
		return nil
	}
	return errors.New("not found article")
}

func (m *TestADBRepo) DeleteGoodArticle(articleID, mainID int) error {
	if articleID == 1 {
		return nil
	}
	return errors.New("not found article")
}
