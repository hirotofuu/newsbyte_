package dbrepo

import (
	"context"
	"database/sql"

	"github.com/hirotofuu/newsbyte/internal/models"
)

type ArticlePostgresDBRepo struct {
	DB *sql.DB
}

func (m *ArticlePostgresDBRepo) InsertArticle(article models.Article) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `insert into articles (title, content, work, main_img, medium, comment_ok, user_id, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`
	var newID int

	err := m.DB.QueryRowContext(ctx, stmt,
		article.Title,
		article.Content,
		article.Work,
		article.MainImg,
		article.Medium,
		article.CommentOK,
		article.UserID,
		article.CreatedAt,
		article.UpdatedAt,
	).Scan(&newID)
	if err != nil {
		return 0, err
	}

	return newID, nil

}

func (m *ArticlePostgresDBRepo) AllArticles() ([]*models.Article, error) {
	ctx, canceal := context.WithTimeout(context.Background(), dbTimeout)
	defer canceal()

	query := `
	select 
		a.id, a.title, a.content, a.work, a.main_img, a.medium, a.comment_ok, a.user_id,  a.created_at, a.updated_at,
		u.user_name, u.avatar_img
	from 
		articles a
		left join users u on (u.id = a.user_id)
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []*models.Article

	for rows.Next() {
		var article models.Article
		err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.Content,
			&article.Work,
			&article.MainImg,
			&article.Medium,
			&article.CommentOK,
			&article.UserID,
			&article.CreatedAt,
			&article.UpdatedAt,
			&article.Name,
			&article.Avatar,
		)
		if err != nil {
			return nil, err
		}

		articles = append(articles, &article)
	}

	return articles, nil
}

func (m *ArticlePostgresDBRepo) UserArticles(userID int) ([]*models.Article, error) {
	ctx, canceal := context.WithTimeout(context.Background(), dbTimeout)
	defer canceal()

	query := `
	select 
		a.id, a.title, a.content, a.work, a.main_img, a.medium, a.comment_ok, a.user_id,  a.created_at, a.updated_at,
		u.user_name, u.avatar_img
	from 
		articles a
		left join users u on (u.id = a.user_id)
		where 
		    a.user_id = $1`

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []*models.Article

	for rows.Next() {
		var article models.Article
		err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.Content,
			&article.Work,
			&article.MainImg,
			&article.Medium,
			&article.CommentOK,
			&article.UserID,
			&article.CreatedAt,
			&article.UpdatedAt,
			&article.Name,
			&article.Avatar,
		)
		if err != nil {
			return nil, err
		}

		articles = append(articles, &article)
	}

	return articles, nil
}

func (m *ArticlePostgresDBRepo) WorkArticles(work string) ([]*models.Article, error) {
	ctx, canceal := context.WithTimeout(context.Background(), dbTimeout)
	defer canceal()

	query := `
	select 
		a.id, a.title, a.content, a.work, a.main_img, a.medium, a.comment_ok, a.user_id,  a.created_at, a.updated_at,
		u.user_name, u.avatar_img
	from 
		articles a
		left join users u on (u.id = a.user_id)
		where 
		    a.work = $1`

	rows, err := m.DB.QueryContext(ctx, query, work)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []*models.Article

	for rows.Next() {
		var article models.Article
		err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.Content,
			&article.Work,
			&article.MainImg,
			&article.Medium,
			&article.CommentOK,
			&article.UserID,
			&article.CreatedAt,
			&article.UpdatedAt,
			&article.Name,
			&article.Avatar,
		)
		if err != nil {
			return nil, err
		}

		articles = append(articles, &article)
	}

	return articles, nil
}

func (m *ArticlePostgresDBRepo) OneArticle(id int) (*models.Article, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
	select 
		a.id, a.title, a.content, a.work, a.main_img, a.medium, a.comment_ok, a.user_id,  a.created_at, a.updated_at,
		u.user_name, u.avatar_img
	from 
		articles a
		left join users u on (u.id = a.user_id)
		where 
		    a.id = $1`

	var article models.Article
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&article.ID,
		&article.Title,
		&article.Content,
		&article.Work,
		&article.MainImg,
		&article.Medium,
		&article.CommentOK,
		&article.UserID,
		&article.CreatedAt,
		&article.UpdatedAt,
		&article.Name,
		&article.Avatar,
	)

	if err != nil {
		return nil, err
	}

	return &article, nil
}

func (m *ArticlePostgresDBRepo) DeleteArticle(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from articles where id = $1`
	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}
