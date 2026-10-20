package repository

import (
	"context"

	"kheun-Lift-backend-go/internal/auth/model"
	"kheun-Lift-backend-go/internal/queries"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateUser(
	ctx context.Context,
	pool *pgxpool.Pool,
	input model.RegisterInput,
	passwordHash string,
) (model.User, error) {
	var user model.User

	err := pool.QueryRow(
		ctx, queries.CreateUser, input.Username, input.Email, passwordHash,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.ProfileImage,
		&user.BodyWeight,
		&user.Height,
		&user.Gender,
		&user.DateOfBirth,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	return user, err
}

func FindLoginUser(
	ctx context.Context,
	pool *pgxpool.Pool,
	email string,
) (model.LoginUser, string, error) {
	var user model.LoginUser
	var passwordHash string

	err := pool.QueryRow(ctx, queries.FindLoginUser, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&passwordHash,
	)

	return user, passwordHash, err
}
