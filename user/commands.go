package user

import "context"

type Commands interface {
	Insert
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
