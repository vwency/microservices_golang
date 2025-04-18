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

type Claims struct {
	UserID string   `json:"user_id"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
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

func (m *JWTManager) GenerateAccessToken(userID string, roles []string) (string, time.Time, error) {
	expirationTime := time.Now().Add(m.accessTokenTTL)
	return m.generateToken(userID, roles, expirationTime)
}

func (m *JWTManager) GenerateRefreshToken(userID string, roles []string) (string, time.Time, error) {
	expirationTime := time.Now().Add(m.refreshTokenTTL)
	return m.generateToken(userID, roles, expirationTime)
}

func (m *JWTManager) generateToken(userID string, roles []string, expirationTime time.Time) (string, time.Time, error) {
	if userID == "" {
		return "", time.Time{}, errors.New("userID cannot be empty")
	}
	if len(roles) == 0 {
		return "", time.Time{}, errors.New("roles cannot be empty")
	}

	claims := &Claims{
		UserID: userID,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.secretKey))
	fmt.Printf("Generating token for user: %s, roles: %v\n", userID, roles)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expirationTime, nil
}

func (m *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	if tokenString == "" {
		return nil, errors.New("token string is empty")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(m.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
