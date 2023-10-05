package dbrepo

import (
	"context"
	"database/sql"

	"github.com/hirotofuu/newsbyte/internal/models"
)

type CommentPostgresDBRepo struct {
	DB *sql.DB
}

func (m *CommentPostgresDBRepo) InsertComment(comment models.Comment) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `insert into comments (comment, article_id, user_id, created_at, updated_at) values ($1, $2, $3, $4, $5) returning id`
	var newID int

	err := m.DB.QueryRowContext(ctx, stmt,
		comment.Comment,
		comment.ArticleID,
		comment.UserID,
		comment.CreatedAt,
		comment.UpdatedAt,
	).Scan(&newID)
	if err != nil {
		return 0, err
	}

	return newID, nil

}

func (m *CommentPostgresDBRepo) ArticleComments(articleID int) ([]*models.Comment, error) {
	ctx, canceal := context.WithTimeout(context.Background(), dbTimeout)
	defer canceal()

	query := `
	select 
		c.id, c.comment, c.user_id, c.article_id, c.created_at, c.updated_at,
		u.user_name, u.avatar_img
	from 
		comments c
		left join users u on (u.id = c.user_id)
		where 
		    c.article_id = $1`

	rows, err := m.DB.QueryContext(ctx, query, articleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.Comment

	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.Comment,
			&comment.UserID,
			&comment.ArticleID,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.Name,
			&comment.Avatar,
		)
		if err != nil {
			return nil, err
		}

		comments = append(comments, &comment)
	}

	return comments, nil
}

func (m *CommentPostgresDBRepo) UserComments(userID int) ([]*models.Comment, error) {
	ctx, canceal := context.WithTimeout(context.Background(), dbTimeout)
	defer canceal()

	query := `
	select 
		c.id, c.comment, c.user_id, c.article_id, c.created_at, c.updated_at,
		u.user_name, u.avatar_img
	from 
		comments c
		left join users u on (u.id = c.user_id)
		where 
		    c.user_id = $1`

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.Comment

	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.Comment,
			&comment.UserID,
			&comment.ArticleID,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.Name,
			&comment.Avatar,
		)
		if err != nil {
			return nil, err
		}

		comments = append(comments, &comment)
	}

	return comments, nil
}

func (m *CommentPostgresDBRepo) DeleteComment(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from comments where id = $1`
	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}
