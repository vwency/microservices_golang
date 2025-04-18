package user_usecase

import "errors"

type UserParams struct {
	Username string
	Email    string
}
type CreateUserParams struct {
	UserParams
	HashedPassword string
	HashedRt       string
	HashedAt       string
	CreatedAt      string
}

type UpdateTokensParams struct {
	UserID             string
	HashedRefreshToken string
	HashedAccessToken  string
}

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

func (p *UpdateTokensParams) Validate() error {
	if p.UserID == "" {
		return errors.New("username cannot be empty")
	}
	if p.HashedRefreshToken == "" {
		return errors.New("refresh token cannot be empty")
	}
	// if p.HashedAccessToken == "" {
	// 	return errors.New("hashed_at cannot be empty")
	// }
	return nil
}
