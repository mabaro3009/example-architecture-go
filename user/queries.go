package user

import (
	"context"
)

type Queries interface {
	Get
}

type Get interface {
	Get(ctx context.Context, id string) (*User, error)
}
