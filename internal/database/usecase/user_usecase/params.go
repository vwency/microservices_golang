package user_usecase

import "errors"

// UserParams contains common user data
type UserParams struct {
	Username string
	Email    string
}

// CreateUserParams contains parameters for creating a user
type CreateUserParams struct {
	UserParams
	HashedPassword string
	HashedRt       string
	HashedAt       string
}

// UpdateTokensParams contains parameters for updating user tokens
type UpdateTokensParams struct {
	Username string
	HashedRt string
	HashedAt string
}

// Validate method for CreateUserParams
func (p *CreateUserParams) Validate() error {
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

// Validate method for UpdateTokensParams
func (p *UpdateTokensParams) Validate() error {
	if p.Username == "" {
		return errors.New("username cannot be empty")
	}
	if p.HashedRt == "" {
		return errors.New("refresh token cannot be empty")
	}
	if p.HashedAt == "" {
		return errors.New("hashed_at cannot be empty")
	}
	return nil
}
