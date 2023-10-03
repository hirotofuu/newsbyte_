package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

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
		{"valid user", `{"user_name":"hiroto","email":"adminexample.com","password":"secret"}`, http.StatusAccepted},
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
