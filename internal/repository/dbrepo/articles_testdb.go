package dbrepo

import (
	"errors"
	"time"

	"github.com/hirotofuu/newsbyte/internal/models"
)

// 記事関連のデモ関数

type TestADBRepo struct{}

// 記事挿入関数
func (m *TestADBRepo) InsertArticle(article models.Article) (int, error) {

	return 1, nil
}

// 記事編集関数
func (m *TestADBRepo) UpdateArticle(article models.Article) error {

	return nil
}

// 全ての記事
func (m *TestADBRepo) AllArticles() ([]*models.Article, error) {
	var articles []*models.Article

	return articles, nil
}

// ユーザーごと記事
func (m *TestADBRepo) UserArticles(userID int) ([]*models.Article, error) {
	var articles []*models.Article

	return articles, nil
}

// 下書き記事
func (m *TestADBRepo) UserSaveArticles(userID int) ([]*models.Article, error) {
	var articles []*models.Article

	return articles, nil
}

// カテゴリごと
func (m *TestADBRepo) WorkArticles(work string) ([]*models.Article, error) {
	var articles []*models.Article

	return articles, nil
}

// 削除
func (m *TestADBRepo) DeleteArticle(id int) error {

	if id == 1 {
		return nil
	}
	return errors.New("not found article")
}

// 複数削除
func (m *TestADBRepo) DeleteSomeArticles(ids []int) error {
	isFlag := true
	for v := range ids {
		if v != 1 {
			isFlag = true
		}
	}

	if isFlag {
		return nil
	}
	return errors.New("not found article")
}

// 一つの記事
func (m *TestADBRepo) OneArticle(id int) (*models.Article, error) {
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

// 編集よう記事fetch
func (m *TestADBRepo) OneEditArticle(id int) (*models.Article, error) {
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

// いいね
func (m *TestADBRepo) InsertGoodArticle(id, mainID int) error {
	if id == 1 {
		return nil
	}
	return errors.New("not found article")
}

// よくないね
func (m *TestADBRepo) DeleteGoodArticle(articleID, mainID int) error {
	if articleID == 1 {
		return nil
	}
	return errors.New("not found article")
}

// いいねしている？
func (m *TestADBRepo) StateGoodArticle(id int) ([]int, error) {
	a := []int{}
	return a, nil
}
