package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	SecretKey       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func NewJWTManager(secretKey string, accessTTL, refreshTTL time.Duration) (*JWTManager, error) {
	return &JWTManager{
		SecretKey:       secretKey,
		AccessTokenTTL:  accessTTL,
		RefreshTokenTTL: refreshTTL,
	}, nil
}

func (m *JWTManager) GenerateAccessToken(payload map[string]interface{}) (string, time.Time, error) {
	expirationTime := time.Now().Add(m.AccessTokenTTL)
	return m.generateToken(payload, expirationTime)
}

func (m *JWTManager) GenerateRefreshToken(payload map[string]interface{}) (string, time.Time, error) {
	expirationTime := time.Now().Add(m.RefreshTokenTTL)
	return m.generateToken(payload, expirationTime)
}

func (m *JWTManager) generateToken(payload map[string]interface{}, expirationTime time.Time) (string, time.Time, error) {
	if len(payload) == 0 {
		return "", time.Time{}, errors.New("payload cannot be empty")
	}

	claims := jwt.MapClaims{}
	for k, v := range payload {
		claims[k] = v
	}
	claims["exp"] = expirationTime.Unix()
	claims["iat"] = time.Now().Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.SecretKey))
	fmt.Printf("Generating token with payload: %v\n", payload)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expirationTime, nil
}

func (m *JWTManager) ValidateToken(tokenString string) (map[string]interface{}, error) {
	if tokenString == "" {
		return nil, errors.New("token string is empty")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(m.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		result := make(map[string]interface{})
		for k, v := range claims {
			result[k] = v
		}
		return result, nil
	}

	return nil, errors.New("invalid token")
}
