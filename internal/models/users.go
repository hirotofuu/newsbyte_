package models

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User describes the data for the User type.
type User struct {
	ID               int       `json:"id"`
	UserName         string    `json:"user_name"`
	Email            string    `json:"email"`
	AvatarImg        string    `json:"avatar_img"`
	Profile          string    `json:"profile"`
	Password         string    `json:"-"`
	FollowingUserIDs []*int    `json:"following_user_ids"`
	Token            string    `json:"token,omitempty"`
	RefreshToken     string    `json:"refresh_token,omitempty"`
	FollowingsCount  int       `json:"followings_count,omitempty"`
	FollowedsCount   int       `json:"followeds_count,omitempty"`
	CreatedAt        time.Time `json:"-"`
	UpdatedAt        time.Time `json:"-"`
}

// PasswordMatches uses Go's bcrypt package to compare a user supplied password
// with the hash we have stored for a given user in the database. If the password
// and hash match, we return true; otherwise, we return false.
func (u *User) PasswordMatches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// invalid password
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
