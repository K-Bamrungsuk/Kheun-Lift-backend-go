package service

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

func hashPassword(password string) string {
	salt := make([]byte, 16)
	rand.Read(salt)

	key := argon2.IDKey([]byte(password), salt, 3, 64*1024, 4, 32)

	return fmt.Sprintf(
		"$argon2id$v=%d$m=65536,t=3,p=4$%s$%s",
		argon2.Version,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(key),
	)
}

func verifyPassword(password, encoded string) (bool, error) {
	invalid := errors.New("invalid or unsupported password hash")
	parts := strings.Split(encoded, "$")

	if len(parts) != 6 || parts[0] != "" ||
		parts[1] != "argon2id" || parts[2] != "v=19" {
		return false, invalid
	}

	var memory, iterations uint32
	var threads uint8

	if _, err := fmt.Sscanf(
		parts[3], "m=%d,t=%d,p=%d",
		&memory, &iterations, &threads,
	); err != nil {
		return false, invalid
	}

	if parts[3] != fmt.Sprintf(
		"m=%d,t=%d,p=%d", memory, iterations, threads,
	) {
		return false, invalid
	}

	if threads < 1 || threads > 4 ||
		iterations < 1 || iterations > 3 ||
		memory < 8*uint32(threads) || memory > 65536 {
		return false, invalid
	}

	salt, err := base64.RawStdEncoding.Strict().DecodeString(parts[4])
	if err != nil || len(salt) != 16 {
		return false, invalid
	}

	expected, err := base64.RawStdEncoding.Strict().DecodeString(parts[5])
	if err != nil || len(expected) != 32 {
		return false, invalid
	}

	actual := argon2.IDKey(
		[]byte(password), salt, iterations, memory, threads, 32,
	)

	return subtle.ConstantTimeCompare(actual, expected) == 1, nil
}
