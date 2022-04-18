package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidUsername       = errors.New("invalid username")
	ErrInvalidRole           = errors.New("invalid role. Valid roles are user and admin")
	ErrUsernameAlreadyExists = errors.New("this username is already in use")
	ErrIDAlreadyExists       = errors.New("this ID is already in use")
)

type PasswordValidator interface {
	Validate(password string) error
}

type PasswordHasher interface {
	Hash(password string) ([]byte, error)
}

type CreatorQueries interface {
	GetByID
	GetByUsername
}

type CreatorCommands interface {
	Insert
}

type Creator struct {
	validator PasswordValidator
	hasher    PasswordHasher
	q         CreatorQueries
	cmd       CreatorCommands
}

func NewCreator(v PasswordValidator, h PasswordHasher, q CreatorQueries, cmd CreatorCommands) *Creator {
	return &Creator{
		validator: v,
		hasher:    h,
		q:         q,
		cmd:       cmd,
	}
}

type CreateParams struct {
	ID       string
	Username string
	Password string
	Role     string
}

func (c *Creator) Create(ctx context.Context, params CreateParams) (*User, error) {
	if err := c.checkCreateParams(ctx, params); err != nil {
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

func (c *Creator) checkCreateParams(ctx context.Context, params CreateParams) error {
	if params.ID != "" {
		_, err := c.q.GetByID(ctx, params.ID)
		if err == nil {
			return ErrIDAlreadyExists
		}
		if err != ErrDoesNotExist {
			return err
		}
	}

	if params.Username == "" {
		return ErrInvalidUsername
	}

	if params.Role != "" && params.Role != RoleUser && params.Role != RoleAdmin {
		return ErrInvalidRole
	}

	_, err := c.q.GetByUsername(ctx, params.Username)
	if err == nil {
		return ErrUsernameAlreadyExists
	}
	if err != ErrDoesNotExist {
		return err
	}

	if err = c.validator.Validate(params.Password); err != nil {
		return fmt.Errorf("invalid password: %w", err)
	}

	return nil
}
