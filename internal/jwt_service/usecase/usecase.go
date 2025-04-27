package usecase_jwt

import (
	"fmt"
	"math"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtUsecase struct{}

func NewJwtUsecase() *JwtUsecase {
	return &JwtUsecase{}
}

func (u *JwtUsecase) GenerateToken(secret string, payload map[string]string, expiresIn time.Duration) (string, int64, error) {
	claims := jwt.MapClaims{}
	for key, value := range payload {
		claims[key] = value
	}

	// Limit maximum token lifetime
	maxExpiresIn := time.Duration(math.MaxInt64/1000000000) * time.Second
	if expiresIn > maxExpiresIn {
		expiresIn = maxExpiresIn
	}

	expirationTime := time.Now().Add(expiresIn)
	claims["exp"] = expirationTime.Unix()

	// Create token with appropriate signing method
	var token *jwt.Token
	if secret == "" {
		token = jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	} else {
		token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	}

	// Sign the token
	var signedToken string
	var err error
	if secret == "" {
		// For unsigned tokens
		signedToken, err = token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	} else {
		// For signed tokens
		signedToken, err = token.SignedString([]byte(secret))
	}

	if err != nil {
		return "", 0, fmt.Errorf("failed to sign token: %w", err)
	}
	return signedToken, expirationTime.Unix(), nil
}

func (u *JwtUsecase) ValidateToken(tokenString string, secret string) (string, map[string]string, int64, error) {
	parseOptions := []jwt.ParserOption{jwt.WithoutClaimsValidation()}

	// Determine allowed signing methods based on whether we have a secret
	var allowedMethods []string
	if secret == "" {
		allowedMethods = []string{"none"}
	} else {
		allowedMethods = []string{"HS256"}
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if secret == "" {
			return jwt.UnsafeAllowNoneSignatureType, nil
		}
		return []byte(secret), nil
	}, append(parseOptions, jwt.WithValidMethods(allowedMethods))...)

	if err != nil {
		return "", nil, 0, fmt.Errorf("token validation failed: %w", err)
	}

	if !token.Valid {
		return "", nil, 0, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", nil, 0, fmt.Errorf("invalid token claims format")
	}

	var expiresAt int64
	if exp, ok := claims["exp"].(float64); ok {
		expiresAt = int64(exp)
	} else if exp, ok := claims["exp"].(int64); ok {
		expiresAt = exp
	}

	userID, _ := claims["username"].(string)
	payload := make(map[string]string)
	for key, value := range claims {
		if strVal, ok := value.(string); ok {
			payload[key] = strVal
		}
	}

	return userID, payload, expiresAt, nil
}
