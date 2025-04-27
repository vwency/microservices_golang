package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secretKey       string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewJWTManager(secretKey string, accessTTL, refreshTTL time.Duration) (*JWTManager, error) {
	if secretKey == "" {
		return nil, errors.New("JWT secret key cannot be empty")
	}
	return &JWTManager{
		secretKey:       secretKey,
		accessTokenTTL:  accessTTL,
		refreshTokenTTL: refreshTTL,
	}, nil
}

func (m *JWTManager) GenerateAccessToken(payload map[string]interface{}) (string, time.Time, error) {
	expirationTime := time.Now().Add(m.accessTokenTTL)
	return m.generateToken(payload, expirationTime)
}

func (m *JWTManager) GenerateRefreshToken(payload map[string]interface{}) (string, time.Time, error) {
	expirationTime := time.Now().Add(m.refreshTokenTTL)
	return m.generateToken(payload, expirationTime)
}

func (m *JWTManager) generateToken(payload map[string]interface{}, expirationTime time.Time) (string, time.Time, error) {
	if len(payload) == 0 {
		return "", time.Time{}, errors.New("payload cannot be empty")
	}

	// Создаем claims с нашим payload и стандартными claims
	claims := jwt.MapClaims{}
	for k, v := range payload {
		claims[k] = v
	}
	claims["exp"] = expirationTime.Unix()
	claims["iat"] = time.Now().Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.secretKey))
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
		return []byte(m.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Конвертируем jwt.MapClaims в обычный map[string]interface{}
		result := make(map[string]interface{})
		for k, v := range claims {
			result[k] = v
		}
		return result, nil
	}

	return nil, errors.New("invalid token")
}
