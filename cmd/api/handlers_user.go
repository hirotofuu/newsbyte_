package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
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

	// generate tokens
	tokens, err := app.auth.CreateTokenPair(user)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	folloingIDs, err := app.DB.GetFollowingUserIDs(user.ID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	user.Token = tokens.Token
	user.FollowingUserIDs = folloingIDs

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

	_, err = app.DB.GetUserIdName(requestPayload.UserName)
	if err != nil && err != sql.ErrNoRows {
		app.errorJSON(w, err)
		return
	}
	if err == nil {
		app.errorJSON(w, errors.New("this id_name is already used"), http.StatusBadRequest)
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
	user.FollowingUserIDs = []*int{}

	ID, err := app.DB.InsertUser(user)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	user.ID = ID

	tokens, err := app.auth.CreateTokenPair(&user)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	user.ID = ID
	user.Token = tokens.Token
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

			// generate tokens
			tokens, err := app.auth.CreateTokenPair(user)
			if err != nil {
				app.errorJSON(w, err)
				return
			}

			user.Token = tokens.Token

			folloingIDs, err := app.DB.GetFollowingUserIDs(user.ID)
			if err != nil {
				app.errorJSON(w, err)
				return
			}

			user.FollowingUserIDs = folloingIDs

			refreshCookie := app.auth.GetRefreshCookie(tokens.RefreshToken)
			http.SetCookie(w, refreshCookie)

			app.writeJSON(w, http.StatusOK, user)

		}
	}
}

func (app *application) refreshToken_next(w http.ResponseWriter, r *http.Request) {
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

			// generate tokens
			tokens, err := app.auth.CreateTokenPair(user)
			if err != nil {
				app.errorJSON(w, err)
				return
			}

			refreshCookie := app.auth.GetRefreshCookie(tokens.RefreshToken)
			http.SetCookie(w, refreshCookie)

			app.writeJSON(w, http.StatusOK, tokens)

		}
	}
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, app.auth.GetExpiredRefreshCookie())
	w.WriteHeader(http.StatusAccepted)
}

func (app *application) getAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := app.DB.AllUsers()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, users)
}

func (app *application) getSearchUsers(w http.ResponseWriter, r *http.Request) {
	key_word := chi.URLParam(r, "key_word")
	users, err := app.DB.SearchUsers(key_word)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, users)
}

func (app *application) getOneUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	user, err := app.DB.OneUser(userID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, user)

}

func (app *application) getOneIdNameUser(w http.ResponseWriter, r *http.Request) {
	id_name := chi.URLParam(r, "id_name")

	user, err := app.DB.OneIdNameUser(id_name)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, user)

}

func (app *application) getFollowingUsers(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	users, err := app.DB.FollowingUsers(userID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, users)
}

func (app *application) getFollowedUsers(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	users, err := app.DB.FollowedUsers(userID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, users)
}

func (app *application) InsertFollow(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	yourID := app.isLogin(w, r)
	if yourID == 0 {
		app.errorJSON(w, errors.New("you are not authenticated"), http.StatusUnauthorized)
		return
	}

	err = app.DB.InsertFollow(id, yourID)
	if err != nil {
		app.errorJSON(w, err)
	}

	resp := JSONResponse{
		Error:   false,
		Message: "insert sefety",
	}

	app.writeJSON(w, http.StatusOK, resp)
}

func (app *application) DeleteFollow(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	yourID := app.isLogin(w, r)
	if yourID == 0 {
		app.errorJSON(w, errors.New("you are not authenticated"), http.StatusUnauthorized)
		return
	}

	err = app.DB.DeleteFollow(id, yourID)
	if err != nil {
		app.errorJSON(w, err)
	}

	resp := JSONResponse{
		Error:   false,
		Message: "delete sefety",
	}

	app.writeJSON(w, http.StatusOK, resp)
}
