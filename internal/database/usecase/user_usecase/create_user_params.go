package user_usecase

import "errors"

type CreateUserParams struct {
	Username       string
	HashedPassword string
	HashedRt       string
	HashedAt       string
	Email          string
}

func (p CreateUserParams) Validate() error {
	if p.Username == "" {
		return errors.New("username cannot be empty")
	}
	if p.HashedPassword == "" {
		return errors.New("password cannot be empty")
	}
	if p.Email == "" {
		return errors.New("email cannot be empty")
	}
	return nil
}
