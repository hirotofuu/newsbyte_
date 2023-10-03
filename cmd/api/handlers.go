package main

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

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

	refreshCookie := app.auth.GetRefreshCookie(tokens.RefreshToken)
	http.SetCookie(w, refreshCookie)

	app.writeJSON(w, http.StatusOK, tokens)
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

	_, err = app.DB.InsertUser(user)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	resp := JSONResponse{
		Error:   false,
		Message: "register succeed",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}
