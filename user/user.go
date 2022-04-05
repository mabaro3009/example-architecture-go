package user

import (
	"errors"
	"time"
)

var (
	ErrDoesNotExist = errors.New("user does not exist")
)

type User struct {
	ID             string
	Username       string
	HashedPassword []byte
	Role           Role
	CreatedAt      time.Time
	DeletedAt      *time.Time
}
