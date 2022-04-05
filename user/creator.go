package user

import (
	"context"
	"errors"
	"github.com/google/uuid"
)

var (
	ErrInvalidUsername = errors.New("invalid username")
	ErrInvalidRole     = errors.New("invalid role. Valid roles are user and admin")
)

type PasswordValidator interface {
	Validate(password string) error
}

type PasswordHasher interface {
	Hash(password string) ([]byte, error)
}

type CreatorQuery interface {
	Insert
}

type Creator struct {
	validator PasswordValidator
	hasher    PasswordHasher
	q         CreatorQuery
}

func NewCreator(v PasswordValidator, h PasswordHasher, q CreatorQuery) *Creator {
	return &Creator{
		validator: v,
		hasher:    h,
		q:         q,
	}
}

type RawCreateParams struct {
	ID       string
	Username string
	Password string
	Role     string
}

func (c *Creator) Create(ctx context.Context, params RawCreateParams) error {
	if params.Username == "" {
		return ErrInvalidUsername
	}

	if params.Role != "" && params.Role != RoleUser && params.Role != RoleAdmin {
		return ErrInvalidRole
	}

	if err := c.validator.Validate(params.Password); err != nil {
		return err
	}

	hashedPassword, err := c.hasher.Hash(params.Password)
	if err != nil {
		return err
	}

	insertParams := &InsertParams{
		ID:             params.ID,
		Username:       params.Username,
		Role:           params.Role,
		HashedPassword: hashedPassword,
	}

	if insertParams.ID == "" {
		insertParams.ID = uuid.NewString()
	}

	if insertParams.Role == "" {
		insertParams.Role = RoleUser
	}

	return c.q.Insert(ctx, insertParams)
}
