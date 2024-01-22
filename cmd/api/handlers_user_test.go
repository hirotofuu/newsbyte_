package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

// ユーザー、認証に関するhandlerテスト
// handler内のデータベースを操作する関数はダミー仕様

func Test_app_authenticate(t *testing.T) {
	var theTests = []struct {
		name               string
		requestBody        string
		expectedStatusCode int
	}{
		{"valid user", `{"email":"admin@example.com","password":"secret"}`, http.StatusOK},
		{"invalid email", `{"email":"admin@xample.com","password":"secret"}`, http.StatusBadRequest},
		{"invalid password", `{"email":"admin@xample.com","password":"secet"}`, http.StatusBadRequest},
	}

	for _, e := range theTests {
		var reader io.Reader
		reader = strings.NewReader(e.requestBody)
		req, _ := http.NewRequest("POST", "/login", reader)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(app.authenticate)

		handler.ServeHTTP(rr, req)

		if e.expectedStatusCode != rr.Code {
			t.Errorf("%s: returned wrong status code: expected %d but got %d", e.name, e.expectedStatusCode, rr.Code)
		}

	}
}

func Test_app_register(t *testing.T) {
	var theTests = []struct {
		name               string
		requestBody        string
		expectedStatusCode int
	}{
		{"valid user register", `{"user_name":"hiroto","email":"adminexample.com","password":"secret"}`, http.StatusAccepted},
		{"valid user", `{"user_name":"hiroto","email":"admin@example.com","password":"secret"}`, http.StatusBadRequest},
		{"valid user", `{"user_name":"","email":"admin@example.com","password":"secret"}`, http.StatusBadRequest},
	}

	for _, e := range theTests {
		var reader io.Reader
		reader = strings.NewReader(e.requestBody)
		req, _ := http.NewRequest("POST", "/register", reader)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(app.register)

		handler.ServeHTTP(rr, req)

		if e.expectedStatusCode != rr.Code {
			t.Errorf("%s: returned wrong status code: expected %d but got %d", e.name, e.expectedStatusCode, rr.Code)
		}

	}
}

func Test_app_UserHandlers(t *testing.T) {

	var tests = []struct {
		name           string
		method         string
		json           string
		paramName      string
		paramID        string
		handler        http.HandlerFunc
		expectedStatus int
	}{
		{"valid following", "GET", "", "id", "1", app.getFollowingUsers, http.StatusOK},
		{"valid followed", "GET", "", "id", "1", app.getFollowingUsers, http.StatusOK},
		{"valid search users", "GET", "", "key_word", "hiroto", app.getSearchUsers, http.StatusOK},
		{"valid one user", "GET", "", "id", "1", app.getOneUser, http.StatusOK},
	}

	for _, e := range tests {
		var req *http.Request
		if e.json == "" {
			req, _ = http.NewRequest(e.method, "/", nil)
		} else {
			req, _ = http.NewRequest(e.method, "/", strings.NewReader(e.json))
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
