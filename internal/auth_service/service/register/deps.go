package register

import (
	"time"

	databasev1 "github.com/vwency/microservices_golang/proto/user_database"
)

type Dependencies struct {
	DBClient   databasev1.DatabaseInitServiceClient
	Logger     interface{ Log(...interface{}) }
	JWTManager interface {
		GenerateAccessToken(payload map[string]interface{}) (string, time.Time, error)
		GenerateRefreshToken(payload map[string]interface{}) (string, time.Time, error)
	}
	TokenPepper string
}
