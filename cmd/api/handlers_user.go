package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/hirotofuu/newsbyte/internal/models"
)

// ログイン
func (app *application) authenticate(w http.ResponseWriter, r *http.Request) {
	// read json payload
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate user against database
	user, err := app.DB.GetUserByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("hello credentials"), http.StatusBadRequest)
		return
	}

	// check password
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("credentials"), http.StatusBadRequest)
		return
	}

	// create a jwt user
	u := jwtUser{
		ID:       user.ID,
		UserName: user.UserName,
	}

	// generate tokens
	tokens, err := app.auth.CreateTokenPair(&u)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	user.Token = tokens.Token
	user.RefreshToken = tokens.RefreshToken

	refreshCookie := app.auth.GetRefreshCookie(tokens.RefreshToken)
	http.SetCookie(w, refreshCookie)

	app.writeJSON(w, http.StatusOK, user)
}

// ユーザー登録
func (app *application) register(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		UserName string `json:"user_name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	_, err = app.DB.GetUserByEmail(requestPayload.Email)
	if err != nil && err != sql.ErrNoRows {
		app.errorJSON(w, err)
		return
	}
	if err == nil {
		app.errorJSON(w, errors.New("this adress is already used"), http.StatusBadRequest)
		return
	}

	var user models.User

	user.UserName = requestPayload.UserName
	user.Email = requestPayload.Email
	user.Password = requestPayload.Password
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.Profile = ""
	user.AvatarImg = "http://s3/0315"

	ID, err := app.DB.InsertUser(user)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	u := jwtUser{
		ID:       ID,
		UserName: user.UserName,
	}

	tokens, err := app.auth.CreateTokenPair(&u)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	user.ID = ID
	user.Token = tokens.Token
	user.RefreshToken = tokens.RefreshToken

	refreshCookie := app.auth.GetRefreshCookie(tokens.RefreshToken)
	http.SetCookie(w, refreshCookie)

	app.writeJSON(w, http.StatusAccepted, user)
}

func (app *application) refreshToken(w http.ResponseWriter, r *http.Request) {
	for _, cookie := range r.Cookies() {
		if cookie.Name == app.auth.CookieName {
			claims := &Claims{}
			refreshToken := cookie.Value

			_, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(app.JWTSecret), nil
			})
			if err != nil {
				app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
				return
			}
			fmt.Println(claims.ExpiresAt)

			userId, err := strconv.Atoi(claims.Subject)
			if err != nil {
				app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
				return
			}

			user, err := app.DB.GetUserByID(userId)
			if err != nil {
				app.errorJSON(w, errors.New("hello credentials"), http.StatusBadRequest)
				return
			}

			// create a jwt user
			u := jwtUser{
				ID:       user.ID,
				UserName: user.UserName,
			}

			// generate tokens
			tokens, err := app.auth.CreateTokenPair(&u)
			if err != nil {
				app.errorJSON(w, err)
				return
			}

			user.Token = tokens.Token
			user.RefreshToken = tokens.RefreshToken

			refreshCookie := app.auth.GetRefreshCookie(tokens.RefreshToken)
			http.SetCookie(w, refreshCookie)

			app.writeJSON(w, http.StatusOK, user)

		}
	}
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, app.auth.GetExpiredRefreshCookie())
	w.WriteHeader(http.StatusAccepted)
}
