package register

import "time"

type tokenData struct {
	hashedPassword     string
	accessToken        string
	accessExpiresAt    time.Time
	refreshToken       string
	hashedAccessToken  string
	hashedRefreshToken string
}

type taskResult struct {
	taskType int
	err      error
}
