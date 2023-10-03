package dbrepo

import (
	"database/sql"
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

func (m *TestDBRepo) InsertUser(user models.User) (int, error) {

	return 2, nil
}

func (m *TestDBRepo) AllUsers() ([]*models.User, error) {
	var users []*models.User

	return users, nil
}
