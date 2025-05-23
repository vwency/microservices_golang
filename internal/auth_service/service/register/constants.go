package register

import (
	"sync"

	"github.com/vwency/microservices_golang/utils/authutils"
)

var (
	ultraFastHashParams = &authutils.Argon2Params{
		Memory:      4 * 1024, // 4 MB - минимум для быстроты
		Iterations:  1,        // Минимальные итерации
		Parallelism: 1,        // Один поток для избежания overhead
		SaltLength:  12,       // Уменьшенная соль для токенов (они временные)
		KeyLength:   24,       // Уменьшенная длина ключа для токенов
	}
)

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
