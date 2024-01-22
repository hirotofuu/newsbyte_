package models

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// ユーザーのタイプ
type User struct {
	ID               int       `json:"id"`
	UserName         string    `json:"user_name"`
	IdName           string    `json:"id_name"`
	Email            string    `json:"email"`
	AvatarImg        string    `json:"avatar_img"`
	Profile          string    `json:"profile"`
	Password         string    `json:"-"`
	FollowingUserIDs []*int    `json:"following_user_ids"`
	Token            string    `json:"token,omitempty"`
	FollowingsCount  int       `json:"followings_count,omitempty"`
	FollowedsCount   int       `json:"followeds_count,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// パスワード照合
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
