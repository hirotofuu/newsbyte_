package dbrepo

import (
	"database/sql"
	"errors"
	"time"

	"github.com/hirotofuu/newsbyte/internal/models"
)

type TestDBRepo struct{}

// GetUserByEmail returns one user by email address
func (m *TestDBRepo) GetUserByEmail(email string) (*models.User, error) {
	if email == "admin@example.com" {
		user := models.User{
			ID:        1,
			UserName:  "hiroto",
			Email:     "admin@example.com",
			Password:  "$2a$14$ajq8Q7fbtFRQvXpdCq7Jcuy.Rx1h/L4J60Otx.gyNLbAYctGMJ9tK",
			Profile:   "",
			AvatarImg: "http:/s3/s",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		return &user, nil
	}

	return nil, sql.ErrNoRows

}

func (m *TestDBRepo) GetUserByID(id int) (*models.User, error) {
	if id == 1 {
		user := models.User{
			ID:        1,
			UserName:  "hiroto",
			Email:     "admin@example.com",
			Password:  "$2a$14$ajq8Q7fbtFRQvXpdCq7Jcuy.Rx1h/L4J60Otx.gyNLbAYctGMJ9tK",
			Profile:   "",
			AvatarImg: "http:/s3/s",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		return &user, nil
	}

	return nil, sql.ErrNoRows

}

func (m *TestDBRepo) GetUserIdName(id_name string) (*models.User, error) {
	if id_name == "futo" {
		user := models.User{
			ID:        1,
			UserName:  "hiroto",
			Email:     "admin@example.com",
			Password:  "$2a$14$ajq8Q7fbtFRQvXpdCq7Jcuy.Rx1h/L4J60Otx.gyNLbAYctGMJ9tK",
			Profile:   "",
			AvatarImg: "http:/s3/s",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		return &user, nil
	}

	return nil, sql.ErrNoRows

}

func (m *TestDBRepo) InsertUser(user models.User) (int, error) {

	return 2, nil
}

func (m *TestDBRepo) AllUsers() ([]*models.User, error) {
	var users []*models.User

	return users, nil
}

func (m *TestDBRepo) SearchUsers(keyWord string) ([]*models.User, error) {
	var users []*models.User

	return users, nil
}

func (m *TestDBRepo) InsertFollow(id, mainID int) error {
	if id == 1 {
		return nil
	}
	return errors.New("not found user")
}

func (m *TestDBRepo) DeleteFollow(id, mainID int) error {
	if id == 1 {
		return nil
	}
	return errors.New("not found user")
}

func (m *TestDBRepo) GetFollowingUserIDs(mainID int) ([]*int, error) {
	var id int
	id = 3
	if mainID == 1 {
		var ids []*int
		ids = append(ids, &id)
		return ids, nil
	}
	return nil, errors.New("not found user")
}

func (m *TestDBRepo) OneUser(id int) (*models.User, error) {
	if id == 1 {
		user := models.User{
			ID:              1,
			UserName:        "hiroto",
			Email:           "admin@example.com",
			Password:        "$2a$14$ajq8Q7fbtFRQvXpdCq7Jcuy.Rx1h/L4J60Otx.gyNLbAYctGMJ9tK",
			Profile:         "",
			AvatarImg:       "http:/s3/s",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			FollowingsCount: 12,
		}
		return &user, nil
	}

	return nil, sql.ErrNoRows
}

func (m *TestDBRepo) OneIdNameUser(id_name string) (*models.User, error) {
	if id_name == "futo" {
		user := models.User{
			ID:              1,
			UserName:        "hiroto",
			Email:           "admin@example.com",
			Password:        "$2a$14$ajq8Q7fbtFRQvXpdCq7Jcuy.Rx1h/L4J60Otx.gyNLbAYctGMJ9tK",
			Profile:         "",
			AvatarImg:       "http:/s3/s",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			FollowingsCount: 12,
		}
		return &user, nil
	}

	return nil, sql.ErrNoRows
}

func (m *TestDBRepo) FollowingUsers(id int) ([]*models.User, error) {
	var users []*models.User

	return users, nil
}

func (m *TestDBRepo) FollowedUsers(id int) ([]*models.User, error) {
	var users []*models.User

	return users, nil
}
