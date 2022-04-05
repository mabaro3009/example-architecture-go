package user

import (
	"context"
)

type Queries interface {
	Get
	Insert
}

type Get interface {
	Get(ctx context.Context, id string) (*User, error)
}

type InsertParams struct {
	ID             string
	Username       string
	HashedPassword []byte
	Role           string
}

type Insert interface {
	Insert(ctx context.Context, params *InsertParams) error
}
