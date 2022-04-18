package hash

import (
	"golang.org/x/crypto/bcrypt"
)

type BCrypt struct {
	cost int
}

func NewBCrypt(cost int) *BCrypt {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}

	return &BCrypt{cost: cost}
}

func (h *BCrypt) Hash(password string) ([]byte, error) {
	pass := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return hashedPassword, nil
}
