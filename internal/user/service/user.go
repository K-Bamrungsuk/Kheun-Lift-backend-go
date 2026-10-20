package service

import (
	"context"
	"time"

	authmodel "kheun-Lift-backend-go/internal/auth/model"
	"kheun-Lift-backend-go/internal/user/model"
	"kheun-Lift-backend-go/internal/user/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

func UpdateUser(
	ctx context.Context,
	pool *pgxpool.Pool,
	userID int,
	input model.UpdateUserInput,
) (authmodel.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return repository.UpdateUser(ctx, pool, userID, input)
}
