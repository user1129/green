package main

import (
	"encoding/json"
	"net/http"
	"testing"

	"greenlight.bcc/internal/assert"
)

func TestRegisterUser(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routesTest())
	defer ts.Close()

	const (
		validName     = "eldos"
		validEmail    = "eldos@gmail.com"
		validPassword = "eldos"
	)

	tests := []struct {
		Topic    string
		Name     string
		Email    string
		Password string
		wantCode int
	}{
		{
			Topic:    "Valid submission",
			Name:     validName,
			Email:    validEmail,
			Password: validPassword,
			wantCode: http.StatusCreated,
		},
		{
			Topic:    "User name is not provided",
			Name:     "",
			Email:    validEmail,
			Password: validPassword,
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			Topic:    "Duplicate email",
			Name:     validName,
			Email:    "eldos@gmail.com",
			Password: validPassword,
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			Topic:    "Test for wrong input",
			Name:     validName,
			Email:    validEmail,
			Password: validPassword,
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Topic, func(t *testing.T) {
			inputData := struct {
				Name     string `json:"name"`
				Email    string `json:"email"`
				Password string `json:"password"`
			}{
				Name:     tt.Name,
				Email:    tt.Email,
				Password: tt.Password,
			}

			b, err := json.Marshal(&inputData)
			if err != nil {
				t.Fatal("wrong input data")
			}
			if tt.Topic == "Test for wrong input" {
				b = append(b, 'a')
			}

			code, _, body := ts.postForm(t, "/v1/users", b)
			t.Log(body)
			assert.Equal(t, code, tt.wantCode)

		})
	}
}

func TestActivateUser(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routesTest())
	defer ts.Close()

	const validToken = "fiorlfkdfiddsfjiovngekwfoe"

	tests := []struct {
		Topic    string
		Token    string
		wantCode int
	}{
		{
			Topic:    "Valid submission",
			Token:    validToken,
			wantCode: http.StatusOK,
		},
		{
			Topic:    "Empty token",
			Token:    "",
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			Topic:    "Invalid token",
			Token:    "aaaaaaaaaaaaaaaaaaaaaaaaaa",
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			Topic:    "Test for wrong input",
			Token:    validToken,
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Topic, func(t *testing.T) {
			inputData := struct {
				TokenPlaintext string `json:"token"`
			}{
				TokenPlaintext: tt.Token,
			}

			b, err := json.Marshal(&inputData)
			if err != nil {
				t.Fatal("wrong input data")
			}
			if tt.Topic == "Test for wrong input" {
				b = append(b, 'a')
			}

			code, _, body := ts.putForm(t, "/v1/users/activated", b)
			t.Log(body)
			assert.Equal(t, code, tt.wantCode)

		})
	}
}
