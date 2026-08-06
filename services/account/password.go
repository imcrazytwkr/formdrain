package account

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Argon2id parameters (see plan_auth.md).
//
// YAGNI: no PasswordHasher interface — one KDF, one encoding, no alternate backends planned.
const (
	argonTime        = 6
	argonMemoryKiB   = 64 * 1024 // 64 MiB
	argonParallelism = 1
	argonSaltLen     = 16
	argonKeyLen      = 32

	maxPasswordBytes = 1024
)

// HashPassword returns a PHC-encoded Argon2id hash.
func HashPassword(password string) (string, error) {
	if len(password) < 1 {
		return "", ErrEmptyPassword
	}

	if len(password) > maxPasswordBytes {
		return "", ErrPasswordTooLong
	}

	salt, err := randomSalt()
	if err != nil {
		return "", err
	}

	key := argon2.IDKey([]byte(password), salt, argonTime, argonMemoryKiB, argonParallelism, argonKeyLen)
	return encodePHC(salt, key, argonTime, argonMemoryKiB, argonParallelism), nil
}

// CheckPassword verifies password against a PHC-encoded Argon2id hash.
func CheckPassword(password, encodedHash string) error {
	if len(password) < 1 {
		return ErrEmptyPassword
	}

	if len(password) > maxPasswordBytes {
		return ErrPasswordTooLong
	}

	salt, key, timeCost, memory, parallelism, err := decodePHC(encodedHash)
	if err != nil {
		return err
	}

	got := argon2.IDKey([]byte(password), salt, timeCost, memory, parallelism, uint32(len(key)))
	if subtle.ConstantTimeCompare(got, key) == 0 {
		return ErrInvalidCredentials
	}

	return nil
}

// HashPasswordArgon2id returns the raw Argon2id key for benchmarking.
func HashPasswordArgon2id(password, salt []byte) []byte {
	return argon2.IDKey(password, salt, argonTime, argonMemoryKiB, argonParallelism, argonKeyLen)
}

func randomSalt() ([]byte, error) {
	salt := make([]byte, argonSaltLen)
	_, err := rand.Read(salt)
	return salt, err
}

func encodePHC(salt, key []byte, timeCost, memory uint32, parallelism uint8) string {
	b64 := base64.RawStdEncoding
	return fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		memory,
		timeCost,
		parallelism,
		b64.EncodeToString(salt),
		b64.EncodeToString(key),
	)
}

func decodePHC(encoded string) (salt, key []byte, timeCost, memory uint32, parallelism uint8, err error) {
	parts := strings.Split(encoded, "$")
	// "", "argon2id", "v=19", "m=...,t=...,p=...", salt, key
	if len(parts) != 6 || parts[1] != "argon2id" {
		return nil, nil, 0, 0, 0, ErrInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil || version != argon2.Version {
		return nil, nil, 0, 0, 0, ErrInvalidHash
	}

	for param := range strings.SplitSeq(parts[3], ",") {
		name, value, ok := strings.Cut(param, "=")
		if !ok {
			return nil, nil, 0, 0, 0, ErrInvalidHash
		}

		n, convErr := strconv.ParseUint(value, 10, 32)
		if convErr != nil {
			return nil, nil, 0, 0, 0, ErrInvalidHash
		}

		switch name {
		case "m":
			memory = uint32(n)
		case "t":
			timeCost = uint32(n)
		case "p":
			if n > 255 {
				return nil, nil, 0, 0, 0, ErrInvalidHash
			}

			parallelism = uint8(n)
		default:
			return nil, nil, 0, 0, 0, ErrInvalidHash
		}
	}

	if memory < 1 || timeCost < 1 || parallelism < 1 {
		return nil, nil, 0, 0, 0, ErrInvalidHash
	}

	b64 := base64.RawStdEncoding
	salt, err = b64.DecodeString(parts[4])
	if err != nil || len(salt) < 1 {
		return nil, nil, 0, 0, 0, ErrInvalidHash
	}

	key, err = b64.DecodeString(parts[5])
	if err != nil || len(key) < 1 {
		return nil, nil, 0, 0, 0, ErrInvalidHash
	}

	return salt, key, timeCost, memory, parallelism, nil
}
