package dbrepo

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/hirotofuu/newsbyte/internal/models"
	"golang.org/x/crypto/bcrypt"
)

const dbTimeout = time.Second * 3

type PostgresDBRepo struct {
	DB *sql.DB
}

func (m *PostgresDBRepo) AllUsers() ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, user_name, password, avatar_img, profile,
	created_at, updated_at
	from users `

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User

	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.UserName,
			&user.Password,
			&user.AvatarImg,
			&user.Profile,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

// GetUserByEmail returns one user by email address
func (m *PostgresDBRepo) GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, user_name, password, avatar_img, profile,
			created_at, updated_at from users where email = $1`

	var user models.User
	row := m.DB.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.UserName,
		&user.Password,
		&user.AvatarImg,
		&user.Profile,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *PostgresDBRepo) GetUserByID(id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, user_name, password, avatar_img, profile,
			created_at, updated_at from users where id = $1`

	var user models.User
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.UserName,
		&user.Password,
		&user.AvatarImg,
		&user.Profile,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *PostgresDBRepo) InsertUser(user models.User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return 0, err
	}

	stmt := `insert into users (user_name, email, avatar_img, profile, password, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7) returning id`
	var newID int

	err = m.DB.QueryRowContext(ctx, stmt,
		user.UserName,
		user.Email,
		user.AvatarImg,
		user.Profile,
		hashedPassword,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&newID)
	if err != nil {
		return 0, nil
	}

	return newID, nil
}
