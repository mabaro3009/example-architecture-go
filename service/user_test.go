package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/mabaro3009/example-architecture-go/user"
	"github.com/stretchr/testify/assert"
)

func TestHandleUserCreate(t *testing.T) {
	id := "abc"
	username := "usr"
	password := "1234"
	role := "user"
	buff, _ := json.Marshal(map[string]string{
		"id":       id,
		"username": username,
		"password": password,
		"role":     role,
	})
	testCases := []struct {
		description string
		creatorErr  error
		expStatus   int
	}{
		{
			description: "invalid role",
			creatorErr:  user.ErrInvalidRole,
			expStatus:   http.StatusBadRequest,
		},
		{
			description: "invalid role",
			creatorErr:  user.ErrInvalidUsername,
			expStatus:   http.StatusBadRequest,
		},
		{
			description: "invalid pass",
			creatorErr:  user.ErrPasswordTooSmall,
			expStatus:   http.StatusBadRequest,
		},
		{
			description: "random err",
			creatorErr:  errors.New("random error"),
			expStatus:   http.StatusInternalServerError,
		},
		{
			description: "success",
			creatorErr:  nil,
			expStatus:   http.StatusCreated,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(buff))
			w := httptest.NewRecorder()
			m := &mockCreator{func(ctx context.Context, params user.CreateParams) (*user.User, error) {
				assert.Equal(t, id, params.ID)
				assert.Equal(t, username, params.Username)
				assert.Equal(t, password, params.Password)
				assert.Equal(t, role, params.Role)

				return &user.User{}, tc.creatorErr
			}}

			handleUserCreate(m)(w, r)

			assert.Equal(t, tc.expStatus, w.Result().StatusCode)
			if tc.expStatus != http.StatusCreated {
				return
			}

			var response map[string]interface{}
			_ = json.NewDecoder(w.Result().Body).Decode(&response)

			expKeys := []string{"id", "username", "role"}
			for _, key := range expKeys {
				_, ok := response[key]
				assert.True(t, ok)
			}
		})
	}
}

func TestHandleUserGet(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		userID := "userID"
		r := httptest.NewRequest(http.MethodPost, "/users/{id}", nil)
		r = mux.SetURLVars(r, map[string]string{"id": userID})
		w := httptest.NewRecorder()

		mock := &mockQuery{func(ctx context.Context, id string) (*user.User, error) {
			assert.Equal(t, userID, id)

			return nil, user.ErrDoesNotExist
		}}

		handleUserGet(mock)(w, r)
		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	})

	t.Run("random error", func(t *testing.T) {
		userID := "userID"
		r := httptest.NewRequest(http.MethodPost, "/users/{id}", nil)
		r = mux.SetURLVars(r, map[string]string{"id": userID})
		w := httptest.NewRecorder()

		mock := &mockQuery{func(ctx context.Context, id string) (*user.User, error) {
			assert.Equal(t, userID, id)

			return nil, errors.New("random error")
		}}

		handleUserGet(mock)(w, r)
		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})

	t.Run("success", func(t *testing.T) {
		userID := "userID"
		r := httptest.NewRequest(http.MethodPost, "/users/{id}", nil)
		r = mux.SetURLVars(r, map[string]string{"id": userID})
		w := httptest.NewRecorder()

		mock := &mockQuery{func(ctx context.Context, id string) (*user.User, error) {
			assert.Equal(t, userID, id)

			return &user.User{}, nil
		}}

		handleUserGet(mock)(w, r)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)

		var response map[string]interface{}
		_ = json.NewDecoder(w.Result().Body).Decode(&response)

		expKeys := []string{"id", "username", "role", "created_at", "deleted_at"}
		for _, key := range expKeys {
			_, ok := response[key]
			assert.True(t, ok)
		}
	})
}

type mockCreator struct {
	create func(ctx context.Context, params user.CreateParams) (*user.User, error)
}

func (m *mockCreator) Create(ctx context.Context, params user.CreateParams) (*user.User, error) {
	return m.create(ctx, params)
}

type mockQuery struct {
	getByID func(ctx context.Context, id string) (*user.User, error)
}

func (m *mockQuery) GetByID(ctx context.Context, id string) (*user.User, error) {
	return m.getByID(ctx, id)
}
