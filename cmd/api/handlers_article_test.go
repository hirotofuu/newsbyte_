package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/hirotofuu/newsbyte/internal/models"
)

func Test_app_articleHandlers(t *testing.T) {
	testUser := models.User{
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

	tokens, _ := app.auth.CreateTokenPair(&testUser)

	testCookie := &http.Cookie{
		Name:     "refresh-token",
		Path:     "/",
		Value:    tokens.RefreshToken,
		Expires:  time.Now().Add(app.auth.RefreshExpiry),
		MaxAge:   int(app.auth.RefreshExpiry.Seconds()),
		SameSite: http.SameSiteStrictMode,
		Domain:   "localhost",
		HttpOnly: true,
		Secure:   true,
	}
	var tests = []struct {
		name           string
		method         string
		json           string
		paramName      string
		paramID        string
		handler        http.HandlerFunc
		addCookie      bool
		Cookie         *http.Cookie
		expectedStatus int
	}{
		{"allArticles", "GET", "", "", "", app.GetAllArticles, false, testCookie, http.StatusOK},

		// fetch one test
		{"oneArticle", "GET", "", "id", "1", app.GetOneArticle, true, testCookie, http.StatusOK},
		{"oneArticle invalid", "GET", "", "id", "100", app.GetOneArticle, true, testCookie, http.StatusBadRequest},
		{"oneArticle bad URL param", "GET", "", "id", "Y", app.GetOneArticle, true, testCookie, http.StatusBadRequest},

		// fetch user's test
		{"userArticle", "GET", "", "user_id", "1", app.GetUserArticles, false, testCookie, http.StatusOK},

		// fetch work test
		{"workArticle", "GET", "", "work", "呪術廻戦", app.GetWorkArticles, false, testCookie, http.StatusOK},

		// delete test
		{"deleteArticle", "DELETE", "", "id", "1", app.DeleteArticle, true, testCookie, http.StatusOK},
		{"deleteArticle invalid", "DELETE", "", "id", "9", app.DeleteArticle, true, testCookie, http.StatusBadRequest},
		{"deleteArticle  bad URL param", "DELETE", "", "id", "Y", app.DeleteArticle, true, testCookie, http.StatusBadRequest},

		// insert article test
		{
			"insertArticle valid",
			"PUT",
			`{"title":"a","content":"it's great.","tags_in":["Spring","Summer","Fall","Winter"],"medium":1,"comment_ok":true,"main_img":"http:sss/sss","user_id":1}`,
			"",
			"",
			app.InsertArticle,
			false,
			testCookie,
			http.StatusOK,
		},
		{
			"insertArticle invalid input",
			"PUT",
			`{"title":2,"content":"it's great.","tags_in":"jujukaisen","medium":"1","comment_ok":true,"main_img":"http:sss/sss","user_id":1}`,
			"",
			"",
			app.InsertArticle,
			false,
			testCookie,
			http.StatusBadRequest,
		},

		{"insertGoodArticle valid", "PUT", "", "id", "1", app.InsertGoodArticle, true, testCookie, http.StatusOK},
		{"insertGoodArticle invalid params", "PUT", "", "id", "2", app.InsertGoodArticle, true, testCookie, http.StatusBadRequest},
		{"insertGoodArticle invalid paramsName", "PUT", "", "ide", "2", app.InsertGoodArticle, true, testCookie, http.StatusBadRequest},
		{"insertGoodArticle not cookie", "PUT", "", "id", "2", app.InsertGoodArticle, false, testCookie, http.StatusUnauthorized},

		{"deleteGoodArticle valid", "PUT", "", "id", "1", app.DeleteGoodArticle, true, testCookie, http.StatusOK},
		{"deleteGoodArticle invalid params", "DELETE", "", "id", "2", app.DeleteGoodArticle, true, testCookie, http.StatusBadRequest},
		{"deleteGoodArticle invalid paramsName", "DELETE", "", "ide", "2", app.DeleteGoodArticle, true, testCookie, http.StatusBadRequest},
		{"deleteGoodArticle valid", "PUT", "", "id", "1", app.DeleteGoodArticle, false, testCookie, http.StatusUnauthorized},
	}

	for _, e := range tests {
		var req *http.Request
		if e.json == "" {
			req, _ = http.NewRequest(e.method, "/", nil)
		} else {
			req, _ = http.NewRequest(e.method, "/", strings.NewReader(e.json))
		}
		if e.addCookie {
			req.AddCookie(e.Cookie)
		}

		if e.paramName != "" {
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add(e.paramName, e.paramID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(e.handler)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatus {
			t.Errorf("%s: wrong status returned; expected %d but got %d", e.name, e.expectedStatus, rr.Code)
		}
	}
}

func Test_app_commentHandlers(t *testing.T) {
	testUser := models.User{
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

	tokens, _ := app.auth.CreateTokenPair(&testUser)

	testCookie := &http.Cookie{
		Name:     "refresh-token",
		Path:     "/",
		Value:    tokens.RefreshToken,
		Expires:  time.Now().Add(app.auth.RefreshExpiry),
		MaxAge:   int(app.auth.RefreshExpiry.Seconds()),
		SameSite: http.SameSiteStrictMode,
		Domain:   "localhost",
		HttpOnly: true,
		Secure:   true,
	}

	var tests = []struct {
		name           string
		method         string
		json           string
		paramName      string
		paramID        string
		handler        http.HandlerFunc
		addCookie      bool
		Cookie         *http.Cookie
		expectedStatus int
	}{
		// fetch user's test
		{"userComment", "GET", "", "user_id", "1", app.GetUserComments, true, testCookie, http.StatusOK},
		{"userComment invalid user_id", "GET", "", "user_id", "3", app.GetUserComments, true, testCookie, http.StatusBadRequest},
		{"userComment invalid paramsName", "GET", "", "user_ide", "3", app.GetUserComments, true, testCookie, http.StatusBadRequest},

		//fetch article's test
		{"userComment", "GET", "", "article_id", "1", app.GetArticleComments, true, testCookie, http.StatusOK},
		{"articleComment invalid article_id", "GET", "", "article_id", "2", app.GetArticleComments, true, testCookie, http.StatusBadRequest},
		{"userComment invalid paramsName", "GET", "", "article_ide", "1", app.GetArticleComments, true, testCookie, http.StatusBadRequest},

		// delete comment
		{"userComment", "DELETE", "", "id", "1", app.DeleteComment, true, testCookie, http.StatusOK},
		{"deleteComment invalid id", "DELETE", "", "id", "2", app.DeleteComment, true, testCookie, http.StatusBadRequest},
		{"deleteComment invalid paramName", "DELETE", "", "ide", "2", app.DeleteComment, true, testCookie, http.StatusBadRequest},

		// good
		{"insertGoodComment valid", "PUT", "", "id", "1", app.InsertGoodComment, true, testCookie, http.StatusOK},
		{"insertGoodComment invalid params", "PUT", "", "id", "2", app.InsertGoodComment, true, testCookie, http.StatusBadRequest},
		{"insertGoodComment invalid paramsName", "PUT", "", "ide", "2", app.InsertGoodComment, true, testCookie, http.StatusBadRequest},
		{"deleteGoodComment valid", "DELETE", "", "id", "1", app.DeleteGoodComment, true, testCookie, http.StatusOK},
		{"deleteGoodComment invalid params", "DELETE", "", "id", "2", app.DeleteGoodComment, true, testCookie, http.StatusBadRequest},
		{"deleteGoodComment invalid paramsName", "DELETE", "", "ide", "2", app.DeleteGoodComment, true, testCookie, http.StatusBadRequest},
	}

	for _, e := range tests {
		var req *http.Request
		if e.json == "" {
			req, _ = http.NewRequest(e.method, "/", nil)
		} else {
			req, _ = http.NewRequest(e.method, "/", strings.NewReader(e.json))
		}

		if e.addCookie {
			req.AddCookie(e.Cookie)
		}

		if e.paramName != "" {
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add(e.paramName, e.paramID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(e.handler)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatus {
			t.Errorf("%s: wrong status returned; expected %d but got %d", e.name, e.expectedStatus, rr.Code)
		}
	}
}
