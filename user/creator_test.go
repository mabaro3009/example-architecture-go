package user

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreate_IDAlreadyExists(t *testing.T) {
	userID := "abc"
	q := &mockCreatorQueries{
		getByID: func(ctx context.Context, id string) (*User, error) {
			assert.Equal(t, userID, id)
			return &User{}, nil
		},
	}

	c := NewCreator(nil, nil, q, nil)

	params := CreateParams{ID: userID}
	u, err := c.Create(context.Background(), params)
	assert.Nil(t, u)
	assert.ErrorIs(t, err, ErrIDAlreadyExists)
}

func TestCreate_UsernameAlreadyExists(t *testing.T) {
	userID := "1"
	usernameExisting := "abc"
	q := &mockCreatorQueries{
		getByID: func(ctx context.Context, id string) (*User, error) {
			assert.Equal(t, userID, id)
			return nil, ErrDoesNotExist
		},
		getByUsername: func(ctx context.Context, username string) (*User, error) {
			assert.Equal(t, usernameExisting, username)
			return &User{}, nil
		},
	}

	c := NewCreator(nil, nil, q, nil)

	params := CreateParams{
		ID:       userID,
		Username: usernameExisting,
		Role:     "user",
	}
	u, err := c.Create(context.Background(), params)
	assert.Nil(t, u)
	assert.ErrorIs(t, err, ErrUsernameAlreadyExists)
}

func TestCreate(t *testing.T) {
	pvError := errors.New("not a good pass")
	testCases := []struct {
		description string
		id          string
		username    string
		password    string
		role        string
		errPV       error
		expError    error
	}{
		{
			description: "invalid username",
			id:          "1",
			username:    "",
			password:    "aa",
			role:        "user",
			errPV:       nil,
			expError:    ErrInvalidUsername,
		},
		{
			description: "invalid role",
			id:          "1",
			username:    "abc",
			password:    "aa",
			role:        "not a role",
			errPV:       nil,
			expError:    ErrInvalidRole,
		},
		{
			description: "invalid password",
			id:          "1",
			username:    "abc",
			password:    "aa",
			role:        "user",
			errPV:       pvError,
			expError:    pvError,
		},
		{
			description: "invalid role",
			id:          "1",
			username:    "abc",
			password:    "aa",
			role:        "not a role",
			errPV:       nil,
			expError:    ErrInvalidRole,
		},
		{
			description: "all good",
			id:          "1",
			username:    "abc",
			password:    "aa",
			role:        "admin",
			errPV:       nil,
			expError:    nil,
		},
		{
			description: "all good with no id",
			id:          "",
			username:    "abc",
			password:    "aa",
			role:        "admin",
			errPV:       nil,
			expError:    nil,
		},
		{
			description: "all good with no role",
			id:          "1",
			username:    "abc",
			password:    "aa",
			role:        "",
			errPV:       nil,
			expError:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {

			v := &mockPassValidator{validate: func(password string) error {
				assert.Equal(t, tc.password, password)

				return tc.errPV
			}}

			h := &mockPassHasher{hash: func(password string) ([]byte, error) {
				assert.Equal(t, tc.password, password)

				return []byte(password), nil
			}}

			q := &mockCreatorQueries{
				getByID: func(ctx context.Context, id string) (*User, error) {
					assert.Equal(t, tc.id, id)

					return nil, ErrDoesNotExist
				},
				getByUsername: func(ctx context.Context, username string) (*User, error) {
					assert.Equal(t, tc.username, username)

					return nil, ErrDoesNotExist
				},
			}

			cmd := &mockCreatorCMD{func(ctx context.Context, params *InsertParams) error {
				assert.Equal(t, tc.username, params.Username)
				assert.Equal(t, []byte(tc.password), params.HashedPassword)

				return nil
			}}

			c := NewCreator(v, h, q, cmd)

			params := CreateParams{
				ID:       tc.id,
				Username: tc.username,
				Password: tc.password,
				Role:     tc.role,
			}
			_, err := c.Create(context.Background(), params)
			assert.ErrorIs(t, err, tc.expError)
		})
	}
}

type mockPassValidator struct {
	validate func(password string) error
}

func (m *mockPassValidator) Validate(password string) error {
	return m.validate(password)
}

type mockPassHasher struct {
	hash func(password string) ([]byte, error)
}

func (m *mockPassHasher) Hash(password string) ([]byte, error) {
	return m.hash(password)
}

type mockCreatorQueries struct {
	getByID       func(ctx context.Context, id string) (*User, error)
	getByUsername func(ctx context.Context, username string) (*User, error)
}

func (m *mockCreatorQueries) GetByID(ctx context.Context, id string) (*User, error) {
	return m.getByID(ctx, id)
}

func (m *mockCreatorQueries) GetByUsername(ctx context.Context, username string) (*User, error) {
	return m.getByUsername(ctx, username)
}

type mockCreatorCMD struct {
	insert func(ctx context.Context, params *InsertParams) error
}

func (m *mockCreatorCMD) Insert(ctx context.Context, params *InsertParams) error {
	return m.insert(ctx, params)
}
