package user

import (
	"context"
)

type Queries interface {
	GetByID
	GetByUsername
}

type GetByID interface {
	GetByID(ctx context.Context, id string) (*User, error)
}

type GetByUsername interface {
	GetByUsername(ctx context.Context, username string) (*User, error)
}
