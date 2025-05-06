package service

type User struct {
	UserID             string
	Username           string
	Email              string
	HashedPassword     string
	HashedRefreshToken string
	HashedAccessToken  string
}
