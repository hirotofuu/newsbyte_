package dbrepo

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"

	"github.com/hirotofuu/newsbyte/internal/models"
)

type ArticlePostgresDBRepo struct {
	DB *sql.DB
}

func (m *ArticlePostgresDBRepo) InsertArticle(article models.Article) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `insert into articles (title, content, tags, medium, comment_ok, is_open_flag, user_id, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`
	var newID int

	err := m.DB.QueryRowContext(ctx, stmt,
		article.Title,
		article.Content,
		article.TagsIn,
		article.Medium,
		article.CommentOK,
		article.IsOpenFlag,
		article.UserID,
		article.CreatedAt,
		article.UpdatedAt,
	).Scan(&newID)
	if err != nil {
		return 0, err
	}

	return newID, nil

}

func (m *ArticlePostgresDBRepo) UpdateArticle(article models.Article) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `update articles set title = $1, content = $2, tags = $3, medium = $4, comment_ok = $5, is_open_flag = $6, updated_at = $7 where id = $8`

	_, err := m.DB.ExecContext(ctx, stmt,
		article.Title,
		article.Content,
		article.TagsIn,
		article.Medium,
		article.CommentOK,
		article.IsOpenFlag,
		article.UpdatedAt,
		article.ID,
	)

	if err != nil {
		return err
	}
	return nil

}

func (m *ArticlePostgresDBRepo) AllArticles() ([]*models.Article, error) {
	ctx, canceal := context.WithTimeout(context.Background(), dbTimeout)
	defer canceal()

	query := `
	select 
		a.id, a.title, a.content, a.tags, a.medium, a.comment_ok, a.user_id,  a.created_at, a.updated_at,
		u.user_name, u.avatar_img, u.id_name
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
			&article.TagsOut,
			&article.Medium,
			&article.CommentOK,
			&article.UserID,
			&article.CreatedAt,
			&article.UpdatedAt,
			&article.Name,
			&article.Avatar,
			&article.IdName,
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
		a.id, a.title, a.content, a.tags, a.medium, a.comment_ok, a.user_id,  a.created_at, a.updated_at,
		u.user_name, u.avatar_img, u.id_name
	from 
		articles a
		left join users u on (u.id = a.user_id)
		where 
		    a.user_id = $1 and
				a.is_open_flag=true
				`

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
			&article.TagsOut,
			&article.Medium,
			&article.CommentOK,
			&article.UserID,
			&article.CreatedAt,
			&article.UpdatedAt,
			&article.Name,
			&article.Avatar,
			&article.IdName,
		)
		if err != nil {
			return nil, err
		}

		articles = append(articles, &article)
	}

	return articles, nil
}

func (m *ArticlePostgresDBRepo) UserSaveArticles(userID int) ([]*models.Article, error) {
	ctx, canceal := context.WithTimeout(context.Background(), dbTimeout)
	defer canceal()

	query := `
	select 
		a.id, a.title, a.content, a.tags, a.medium, a.comment_ok, a.user_id,  a.created_at, a.updated_at,
		u.user_name, u.avatar_img
	from 
		articles a
		left join users u on (u.id = a.user_id)
		where 
		    a.user_id = $1 and
				a.is_open_flag=false
				`

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
			&article.TagsOut,
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
		a.id, a.title, a.tags, a.content, a.medium, a.comment_ok, a.user_id,  a.created_at, a.updated_at,
		u.user_name, u.avatar_img, u.id_name
	from 
		articles a
		left join users u on (u.id = a.user_id)
		where 
		a.is_open_flag=true`

	reg := `( |　)+`
	keyArr1 := regexp.MustCompile(reg).Split(work, -1)
	for _, s := range keyArr1 {
		query = query + fmt.Sprintf(" and '%s' = ANY(a.tags)", s)
	}

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
			&article.TagsOut,
			&article.Content,
			&article.Medium,
			&article.CommentOK,
			&article.UserID,
			&article.CreatedAt,
			&article.UpdatedAt,
			&article.Name,
			&article.Avatar,
			&article.IdName,
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
		a.id, a.title, a.content, a.tags, a.medium, a.comment_ok, a.user_id, a.created_at, a.updated_at, u.user_name, u.avatar_img, u.id_name
	from 
		articles a
		left join users u on (u.id = a.user_id)
		where 
		    a.id = $1 and
				a.is_open_flag=true`

	var article models.Article
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&article.ID,
		&article.Title,
		&article.Content,
		&article.TagsOut,
		&article.Medium,
		&article.CommentOK,
		&article.UserID,
		&article.CreatedAt,
		&article.UpdatedAt,
		&article.Name,
		&article.Avatar,
		&article.IdName,
	)

	if err != nil {
		return nil, err
	}

	return &article, nil
}

func (m *ArticlePostgresDBRepo) OneEditArticle(id int) (*models.Article, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
	select 
		a.id, a.title, a.content, a.tags, a.medium, a.comment_ok, a.user_id, a.created_at, a.updated_at, u.user_name, u.avatar_img, u.id_name, a.is_open_flag
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
		&article.TagsOut,
		&article.Medium,
		&article.CommentOK,
		&article.UserID,
		&article.CreatedAt,
		&article.UpdatedAt,
		&article.Name,
		&article.Avatar,
		&article.IdName,
		&article.IsOpenFlag,
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

func (m *ArticlePostgresDBRepo) DeleteSomeArticles(ids []int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from articles where id in ( `
	for i, v := range ids {
		if i == 0 {
			stmt = stmt + fmt.Sprintf("%d", v)
		} else {
			stmt = stmt + fmt.Sprintf(",%d", v)
		}
	}
	stmt = stmt + ")"
	
	_, err := m.DB.ExecContext(ctx, stmt)
	if err != nil {
		return err
	}

	return nil
}

func (m *ArticlePostgresDBRepo) InsertGoodArticle(id, mainID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `insert into article_goods (article_id, user_id) values ($1, $2) returning id`
	var newID int

	err := m.DB.QueryRowContext(ctx, stmt, id, mainID).Scan(&newID)
	if err != nil {
		return err
	}

	return nil

}

func (m *ArticlePostgresDBRepo) StateGoodArticle(id int) ([]int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select user_id from article_goods where article_id = $1`
	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIDs []int

	for rows.Next() {
		var userID int
		err := rows.Scan(&userID)
		if err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}
	if err != nil {
		return nil, err
	}

	return userIDs, nil

}

func (m *ArticlePostgresDBRepo) DeleteGoodArticle(articleID, mainID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from article_goods where article_id = $1 and user_id = $2`

	_, err := m.DB.ExecContext(ctx, stmt, articleID, mainID)
	if err != nil {
		return err
	}

	return nil

}
