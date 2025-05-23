package authutils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"unsafe"

	"golang.org/x/crypto/argon2"
)

type Argon2Params struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

var (
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

var DefaultArgon2Params = &Argon2Params{
	Memory:      16 * 1024, // 16 MB
	Iterations:  3,
	Parallelism: 2,
	SaltLength:  16,
	KeyLength:   32,
}

// Пулы для переиспользования буферов
var (
	saltPool = sync.Pool{
		New: func() interface{} {
			return make([]byte, 32) // Максимальный размер соли
		},
	}

	passwordPool = sync.Pool{
		New: func() interface{} {
			return make([]byte, 0, 256) // Предполагаемый максимальный размер пароля + userID
		},
	}

	builderPool = sync.Pool{
		New: func() interface{} {
			return &strings.Builder{}
		},
	}
)

// Pre-allocated strings для избежания аллокаций
const (
	argon2Prefix  = "$argon2id$v=19$m="
	hashSeparator = "$"
	commaT        = ",t="
	commaP        = ",p="
	mPrefix       = "m="
	tPrefix       = "t="
	pPrefix       = "p="
	vPrefix       = "v=19"
)

// Оптимизированная функция для объединения строк без аллокаций
func unsafeStringToBytes(s string) []byte {
	return (*[0x7fff0000]byte)(unsafe.Pointer(
		(*(*uintptr)(unsafe.Pointer(&s)))),
	)[:len(s):len(s)]
}

func GenHash(userID, password string, p *Argon2Params) (string, error) {
	if p == nil {
		p = DefaultArgon2Params
	}

	// Получаем буфер соли из пула
	saltBuf := saltPool.Get().([]byte)
	defer saltPool.Put(saltBuf)

	salt := saltBuf[:p.SaltLength]
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// Получаем буфер для пароля из пула
	passwordBuf := passwordPool.Get().([]byte)
	defer func() {
		// Очищаем конфиденциальные данные перед возвратом в пул
		for i := range passwordBuf {
			passwordBuf[i] = 0
		}
		passwordPool.Put(passwordBuf[:0])
	}()

	// Создаем peppered password без дополнительных аллокаций
	passwordBuf = passwordBuf[:0]
	passwordBuf = append(passwordBuf, password...)
	passwordBuf = append(passwordBuf, userID...)

	hash := argon2.IDKey(
		passwordBuf,
		salt,
		p.Iterations,
		p.Memory,
		p.Parallelism,
		p.KeyLength,
	)

	// Получаем builder из пула
	builder := builderPool.Get().(*strings.Builder)
	defer func() {
		builder.Reset()
		builderPool.Put(builder)
	}()

	// Предварительно рассчитываем размер для минимизации реаллокаций
	estimatedSize := len(argon2Prefix) + 20 + len(commaT) + 10 + len(commaP) + 5 +
		len(hashSeparator)*2 + base64.RawStdEncoding.EncodedLen(int(p.SaltLength)) +
		base64.RawStdEncoding.EncodedLen(int(p.KeyLength))
	builder.Grow(estimatedSize)

	// Используем более эффективное форматирование
	builder.WriteString(argon2Prefix)
	writeUint32(builder, p.Memory)
	builder.WriteString(commaT)
	writeUint32(builder, p.Iterations)
	builder.WriteString(commaP)
	builder.WriteByte(byte('0' + p.Parallelism)) // Более быстрое преобразование для малых чисел
	builder.WriteString(hashSeparator)

	// Кодируем напрямую в builder
	writeBBase64(builder, salt)
	builder.WriteString(hashSeparator)
	writeBBase64(builder, hash)

	return builder.String(), nil
}

// Оптимизированная функция записи uint32
func writeUint32(b *strings.Builder, n uint32) {
	if n < 10 {
		b.WriteByte(byte('0' + n))
	} else {
		b.WriteString(strconv.FormatUint(uint64(n), 10))
	}
}

// Оптимизированная функция записи base64 напрямую в builder
func writeBBase64(b *strings.Builder, data []byte) {
	encodedLen := base64.RawStdEncoding.EncodedLen(len(data))
	oldLen := b.Len()

	// Расширяем builder на нужный размер
	for i := 0; i < encodedLen; i++ {
		b.WriteByte(0)
	}

	// Получаем slice для записи
	str := (*b).String()
	buf := unsafeStringToBytes(str)[oldLen : oldLen+encodedLen]
	base64.RawStdEncoding.Encode(buf, data)
}

func ComparePasswordAndHash(key, password, encodedHash string) (bool, error) {
	p, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Получаем буфер из пула
	passwordBuf := passwordPool.Get().([]byte)
	defer func() {
		// Очищаем конфиденциальные данные
		for i := range passwordBuf {
			passwordBuf[i] = 0
		}
		passwordPool.Put(passwordBuf[:0])
	}()

	// Создаем peppered password
	passwordBuf = passwordBuf[:0]
	passwordBuf = append(passwordBuf, password...)
	passwordBuf = append(passwordBuf, key...)

	otherHash := argon2.IDKey(passwordBuf, salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

	return subtle.ConstantTimeCompare(hash, otherHash) == 1, nil
}

// Кеш для часто используемых параметров
var paramsCache = sync.Map{}

type cacheKey struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
}

func decodeHash(encodedHash string) (*Argon2Params, []byte, []byte, error) {
	// Быстрая проверка формата без Split
	if len(encodedHash) < 30 || !strings.HasPrefix(encodedHash, "$argon2id$") {
		return nil, nil, nil, ErrInvalidHash
	}

	vals := strings.Split(encodedHash, hashSeparator)
	if len(vals) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	if vals[2] != vPrefix {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	// Оптимизированный парсинг параметров
	paramsPart := vals[3]

	// Находим позиции разделителей
	tPos := strings.Index(paramsPart, commaT)
	pPos := strings.Index(paramsPart, commaP)

	if tPos == -1 || pPos == -1 || tPos >= pPos {
		return nil, nil, nil, ErrInvalidHash
	}

	// Парсим параметры
	memStr := paramsPart[2:tPos] // Пропускаем "m="
	mem, err := parseUint32Fast(memStr)
	if err != nil {
		return nil, nil, nil, ErrInvalidHash
	}

	iterStr := paramsPart[tPos+3 : pPos] // Пропускаем ",t="
	iter, err := parseUint32Fast(iterStr)
	if err != nil {
		return nil, nil, nil, ErrInvalidHash
	}

	parallelStr := paramsPart[pPos+3:] // Пропускаем ",p="
	parallel, err := strconv.ParseUint(parallelStr, 10, 8)
	if err != nil {
		return nil, nil, nil, ErrInvalidHash
	}

	// Проверяем кеш для параметров
	key := cacheKey{mem, iter, uint8(parallel)}
	if cached, ok := paramsCache.Load(key); ok {
		p := cached.(*Argon2Params)

		salt, err := base64.RawStdEncoding.DecodeString(vals[4])
		if err != nil {
			return nil, nil, nil, err
		}

		hash, err := base64.RawStdEncoding.DecodeString(vals[5])
		if err != nil {
			return nil, nil, nil, err
		}

		return p, salt, hash, nil
	}

	// Создаем новые параметры
	salt, err := base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}

	hashBytes, err := base64.RawStdEncoding.DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}

	p := &Argon2Params{
		Memory:      mem,
		Iterations:  iter,
		Parallelism: uint8(parallel),
		SaltLength:  uint32(len(salt)),
		KeyLength:   uint32(len(hashBytes)),
	}

	// Сохраняем в кеш
	paramsCache.Store(key, p)

	return p, salt, hashBytes, nil
}

// Быстрый парсинг uint32 для небольших чисел
func parseUint32Fast(s string) (uint32, error) {
	if len(s) == 0 {
		return 0, fmt.Errorf("empty string")
	}

	// Для чисел до 4 цифр используем быстрый алгоритм
	if len(s) <= 4 {
		var result uint32
		for i, c := range []byte(s) {
			if c < '0' || c > '9' {
				return 0, fmt.Errorf("invalid character")
			}
			digit := uint32(c - '0')
			switch len(s) - i - 1 {
			case 0:
				result += digit
			case 1:
				result += digit * 10
			case 2:
				result += digit * 100
			case 3:
				result += digit * 1000
			}
		}
		return result, nil
	}

	// Для больших чисел используем стандартный парсинг
	val, err := strconv.ParseUint(s, 10, 32)
	return uint32(val), err
}
