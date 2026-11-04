package repository

import (
	"context"

	authmodel "kheun-Lift-backend-go/internal/auth/model"
	"kheun-Lift-backend-go/internal/queries"
	"kheun-Lift-backend-go/internal/user/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

func UpdateUser(
	ctx context.Context,
	pool *pgxpool.Pool,
	userID int,
	input model.UpdateUserInput,
) (authmodel.User, error) {
	var user authmodel.User

	err := pool.QueryRow(
		ctx, queries.UpdateUser, userID,
		input.Username, input.BodyWeight, input.Height,
		input.Gender, input.DateOfBirth, input.ProfileImage,
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
