package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"log"
	"time"

	"kheun-Lift-backend-go/internal/auth/model"
	"kheun-Lift-backend-go/internal/auth/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrInvalidCredentials = errors.New("invalid email or password")

func RegisterUser(
	ctx context.Context,
	pool *pgxpool.Pool,
	input model.RegisterInput,
) (model.User, error) {
	passwordHash := hashPassword(input.Password)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return repository.CreateUser(ctx, pool, input, passwordHash)
}

func Login(
	ctx context.Context,
	pool *pgxpool.Pool,
	input model.LoginInput,
) (model.LoginResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	user, passwordHash, err := repository.FindLoginUser(ctx, pool, input.Email)
	missingUser := errors.Is(err, pgx.ErrNoRows)

	if err != nil && !missingUser {
		log.Print("Find login user failed")
		return model.LoginResult{}, err
	}

	if missingUser {
		passwordHash = "$argon2id$v=19$m=65536,t=3,p=4$" +
			"AAAAAAAAAAAAAAAAAAAAAA$" +
			"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
	}

	matches, err := verifyPassword(input.Password, passwordHash)
	if err != nil {
		log.Print("Verify password hash failed")
		return model.LoginResult{}, err
	}
	if !matches || missingUser {
		return model.LoginResult{}, ErrInvalidCredentials
	}

	token := make([]byte, 32)
	rand.Read(token)

	value := base64.RawURLEncoding.EncodeToString(token)
	tokenHash := sha256.Sum256([]byte(value))
	expiresAt := time.Now().UTC().Add(24 * time.Hour)

	if err := repository.CreateSession(
		ctx, pool, tokenHash, user.ID, expiresAt,
	); err != nil {
		log.Print("Create session failed")
		return model.LoginResult{}, err
	}

	return model.LoginResult{
		User:      user,
		Token:     value,
		ExpiresAt: expiresAt,
	}, nil
}

func FindSession(
	ctx context.Context,
	pool *pgxpool.Pool,
	token string,
) (model.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tokenHash := sha256.Sum256([]byte(token))
	return repository.FindSession(ctx, pool, tokenHash)
}

func Logout(
	ctx context.Context,
	pool *pgxpool.Pool,
	tokenHash [32]byte,
) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return repository.DeleteSession(ctx, pool, tokenHash)
}
