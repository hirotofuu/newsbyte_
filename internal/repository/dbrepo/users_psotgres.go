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

	query := `select id, user_name, id_name, avatar_img, profile,
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
			&user.UserName,
			&user.IdName,
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

	query := `select id, user_name, id_name, password, avatar_img, profile,
			created_at, updated_at from users where email = $1`

	var user models.User
	row := m.DB.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.UserName,
		&user.IdName,
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

func (m *PostgresDBRepo) GetUserIdName(id_name string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, user_name, id_name, avatar_img, profile,
			created_at, updated_at from users where id_name = $1`

	var user models.User
	row := m.DB.QueryRowContext(ctx, query, id_name)

	err := row.Scan(
		&user.ID,
		&user.UserName,
		&user.IdName,
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

	query := `select id, email, user_name, id_name, password, avatar_img, profile,
			created_at, updated_at from users where id = $1`

	var user models.User
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.UserName,
		&user.IdName,
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

	stmt := `insert into users (user_name, id_name, email, avatar_img, profile, password, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7, $8) returning id`
	var newID int

	err = m.DB.QueryRowContext(ctx, stmt,
		user.UserName,
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

func (m *PostgresDBRepo) UpdateUser(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `update users set user_name = $1, profile = $2, updated_at = $3 where id = $4`

	_, err := m.DB.ExecContext(ctx, stmt,
		user.UserName,
		user.Profile,
		user.UpdatedAt,
		user.ID,
	)

	if err != nil {
		return err
	}
	return nil

}

func (m *PostgresDBRepo) InsertFollow(id, mainID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `insert into follows (followed_id, following_id) values ($1, $2) returning id`
	var newID int

	err := m.DB.QueryRowContext(ctx, stmt, id, mainID).Scan(&newID)
	if err != nil {
		return err
	}

	return nil

}

func (m *PostgresDBRepo) DeleteFollow(id, mainID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from follows where followed_id = $1 and following_id = $2`

	_, err := m.DB.ExecContext(ctx, stmt, id, mainID)
	if err != nil {
		return err
	}

	return nil

}

func (m *PostgresDBRepo) GetFollowingUserIDs(mainID int) ([]*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `select followed_id from follows f where f.following_id = $1`

	rows, err := m.DB.QueryContext(ctx, query, mainID)
	if err != nil {
		return nil, err
	}

	var followingUserIDs []*int
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}

		followingUserIDs = append(followingUserIDs, &id)
	}

	return followingUserIDs, nil

}

func (m *PostgresDBRepo) SearchUsers(keyWord string) ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, user_name, id_name, avatar_img, profile,
	created_at, updated_at
	from users where user_name = $1`

	rows, err := m.DB.QueryContext(ctx, query, keyWord)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User

	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.UserName,
			&user.IdName,
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

func (m *PostgresDBRepo) OneUser(id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
	select
			u.id, u.user_name, u.id_name, u.avatar_img, u.profile,
			u.created_at, u.updated_at, coalesce(ing.followings_count, 0), coalesce(ed.followeds_count, 0)
	from users u
	left join
		(select count(*)  followings_count, following_id
		from
			follows
		group by following_id)  ing 
	on (ing.following_id = u.id)

	left join
		(select count(*)  followeds_count, followed_id
		from
			follows
		group by followed_id)  ed 
	on (ed.followed_id = u.id)

	where id = $1`

	var user models.User
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.UserName,
		&user.IdName,
		&user.AvatarImg,
		&user.Profile,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.FollowingsCount,
		&user.FollowedsCount,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *PostgresDBRepo) OneIdNameUser(id_name string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
	select
			u.id, u.user_name, u.id_name, u.avatar_img, u.profile,
			u.created_at, u.updated_at, coalesce(ing.followings_count, 0), coalesce(ed.followeds_count, 0)
	from users u
	left join
		(select count(*)  followings_count, following_id
		from
			follows
		group by following_id)  ing 
	on (ing.following_id = u.id)

	left join
		(select count(*)  followeds_count, followed_id
		from
			follows
		group by followed_id)  ed 
	on (ed.followed_id = u.id)

	where id_name = $1`

	var user models.User
	row := m.DB.QueryRowContext(ctx, query, id_name)

	err := row.Scan(
		&user.ID,
		&user.UserName,
		&user.IdName,
		&user.AvatarImg,
		&user.Profile,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.FollowingsCount,
		&user.FollowedsCount,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *PostgresDBRepo) FollowingUsers(id int) ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select u.id, u.user_name, u.id_name, u.avatar_img, u.profile,
	u.created_at, u.updated_at
	from follows f
	left join users u on (f.followed_id= u.id)
	where f.following_id = $1
	`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User

	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.UserName,
			&user.IdName,
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

func (m *PostgresDBRepo) FollowedUsers(id int) ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select u.id, u.user_name, u.id_name, u.avatar_img, u.profile,
	u.created_at, u.updated_at
	from follows f
	left join users u on (f.following_id= u.id)
	where f.followed_id = $1
	`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User

	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.UserName,
			&user.IdName,
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
