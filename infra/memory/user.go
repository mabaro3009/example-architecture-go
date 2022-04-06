package memory

import (
	"context"
	"github.com/mabaro3009/example-architecture-go/user"
	"time"
)

type userMem struct {
	ID             string
	Username       string
	HashedPassword []byte
	Role           string
	CreatedAt      time.Time
	DeletedAt      *time.Time
}

func (u *userMem) ToDomain() *user.User {
	return &user.User{
		ID:             u.ID,
		Username:       u.Username,
		HashedPassword: u.HashedPassword,
		Role:           user.Role(u.Role),
		CreatedAt:      u.CreatedAt,
		DeletedAt:      u.DeletedAt,
	}
}

type UserDB struct {
	users map[string]*userMem
}

func NewUserDB() *UserDB {
	return &UserDB{
		users: make(map[string]*userMem),
	}
}

func (m *UserDB) Insert(_ context.Context, params *user.InsertParams) error {
	u := &userMem{
		ID:             params.ID,
		Username:       params.Username,
		HashedPassword: params.HashedPassword,
		Role:           params.Role,
		CreatedAt:      time.Now(),
		DeletedAt:      nil,
	}

	m.users[params.ID] = u

	return nil
}

func (m *UserDB) GetByID(_ context.Context, id string) (*user.User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, user.ErrDoesNotExist
	}

	return u.ToDomain(), nil
}

func (m *UserDB) GetByUsername(_ context.Context, username string) (*user.User, error) {
	for _, u := range m.users {
		if u.Username == username {
			return u.ToDomain(), nil
		}
	}

	return nil, user.ErrDoesNotExist
}
