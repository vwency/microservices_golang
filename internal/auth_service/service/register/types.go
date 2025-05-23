package register

import (
	"sync"
	"time"

	"github.com/vwency/microservices_golang/utils/authutils"
)

// Константы для типов задач
const (
	taskPassword = iota
	taskAccessToken
	taskRefreshToken
	taskHashAccess
	taskHashRefresh
)

// Параметры для быстрого хэширования токенов
var ultraFastHashParams = &authutils.Argon2Params{
	Memory:      4 * 1024, // 4 MB - минимум для быстроты
	Iterations:  1,        // Минимальные итерации
	Parallelism: 1,        // Один поток для избежания overhead
	SaltLength:  12,       // Уменьшенная соль для токенов (они временные)
	KeyLength:   24,       // Уменьшенная длина ключа для токенов
}

// tokenData содержит все данные, генерируемые в процессе регистрации
type tokenData struct {
	hashedPassword     string
	accessToken        string
	refreshToken       string
	hashedAccessToken  string
	hashedRefreshToken string
	accessExpiresAt    time.Time
}

// reset очищает tokenData для повторного использования из пула
func (td *tokenData) reset() {
	td.hashedPassword = ""
	td.accessToken = ""
	td.refreshToken = ""
	td.hashedAccessToken = ""
	td.hashedRefreshToken = ""
	td.accessExpiresAt = time.Time{}
}

// taskResult представляет результат выполнения задачи
type taskResult struct {
	taskType int
	err      error
}

// Пулы объектов для оптимизации памяти
var (
	tokenDataPool = sync.Pool{
		New: func() interface{} {
			return &tokenData{}
		},
	}

	wgPool = sync.Pool{
		New: func() interface{} {
			return &sync.WaitGroup{}
		},
	}
)

type Logger interface {
	Log(keyvals ...interface{}) error
}
