package user

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"time"
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

type CreatorCommands interface {
	Insert
}

type Creator struct {
	validator PasswordValidator
	hasher    PasswordHasher
	cmd       CreatorCommands
}

func NewCreator(v PasswordValidator, h PasswordHasher, cmd CreatorCommands) *Creator {
	return &Creator{
		validator: v,
		hasher:    h,
		cmd:       cmd,
	}
}

type RawCreateParams struct {
	ID       string
	Username string
	Password string
	Role     string
}

func (c *Creator) Create(ctx context.Context, params RawCreateParams) (*User, error) {
	if params.Username == "" {
		return nil, ErrInvalidUsername
	}

	if params.Role != "" && params.Role != RoleUser && params.Role != RoleAdmin {
		return nil, ErrInvalidRole
	}

	if err := c.validator.Validate(params.Password); err != nil {
		return nil, err
	}

	hashedPassword, err := c.hasher.Hash(params.Password)
	if err != nil {
		return nil, err
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

	if err = c.cmd.Insert(ctx, insertParams); err != nil {
		return nil, err
	}

	return &User{
		ID:             insertParams.ID,
		Username:       insertParams.Username,
		HashedPassword: insertParams.HashedPassword,
		Role:           Role(insertParams.Role),
		CreatedAt:      time.Now(),
		DeletedAt:      nil,
	}, nil
}
