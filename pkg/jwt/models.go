package jwt

type TokenResponse struct {
	AccessToken  string                 `json:"access_token"`
	RefreshToken string                 `json:"refresh_token"`
	ExpiresAt    int64                  `json:"expires_at"`
	Payload      map[string]interface{} `json:"payload,omitempty"`
}

type ValidateRequest struct {
	Token string `json:"token"`
}

type ValidateResponse struct {
	Valid     bool                   `json:"valid"`
	Payload   map[string]interface{} `json:"payload"`
	ExpiresAt int64                  `json:"expires_at"`
}
