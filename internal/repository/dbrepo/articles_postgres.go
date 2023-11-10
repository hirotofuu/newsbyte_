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

func (m *ArticlePostgresDBRepo) AllArticles() ([]*models.Article, error) {
	ctx, canceal := context.WithTimeout(context.Background(), dbTimeout)
	defer canceal()

	query := `
	select 
		a.id, a.title, a.content, a.tags, a.medium, a.comment_ok, a.user_id,  a.created_at, a.updated_at,
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

func (m *ArticlePostgresDBRepo) UserArticles(userID int) ([]*models.Article, error) {
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
		u.user_name, u.avatar_img
	from 
		articles a
		left join users u on (u.id = a.user_id)
		where 
		$1 = ANY(a.tags) and
		a.is_open_flag=true`

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
			&article.TagsOut,
			&article.Content,
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

func (m *ArticlePostgresDBRepo) OneArticle(id, mainID int) (*models.Article, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
	select 
		a.id, a.title, a.content, a.tags, a.medium, a.comment_ok, a.user_id,  a.created_at, a.updated_at, u.user_name, u.avatar_img, coalesce(is_good_flag, 0), coalesce(g.goods_count, 0)
	from 
		articles a
		left join users u on (u.id = a.user_id)

		left join
			(select article_id,
				(case
					when u.id = $2 then 1
					else 0	
				end) is_good_flag
			from article_goods n
			left join
				users u on (u.id = n.user_id)
			group by article_id, u.id
			) m	
		on (a.id = m.article_id)

    left join
      (select count(*) as goods_count, article_id
      from
       article_goods
      group by article_id)  g 
    on (g.article_id = a.id)
		where 
		    a.id = $1 and
				a.is_open_flag=true`

	var article models.Article
	row := m.DB.QueryRowContext(ctx, query, id, mainID)

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
		&article.IsGoodFlag,
		&article.GoodsCount,
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
