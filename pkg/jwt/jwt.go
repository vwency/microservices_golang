package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secretKey       []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewJWTManager(secretKey string, accessTTL, refreshTTL time.Duration) (*JWTManager, error) {
	if secretKey == "" {
		return nil, errors.New("secret key cannot be empty")
	}
	return &JWTManager{
		secretKey:       []byte(secretKey),
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

	claims := make(jwt.MapClaims, len(payload)+2)
	for k, v := range payload {
		claims[k] = v
	}
	claims["exp"] = expirationTime.Unix()
	claims["iat"] = time.Now().Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(m.secretKey)
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
		return m.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		result := make(map[string]interface{}, len(claims))
		for k, v := range claims {
			result[k] = v
		}
		return result, nil
	}

	return nil, errors.New("invalid token")
}
