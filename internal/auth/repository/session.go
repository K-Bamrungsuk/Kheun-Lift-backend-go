package repository

import (
	"context"
	"time"

	"kheun-Lift-backend-go/internal/auth/model"
	"kheun-Lift-backend-go/internal/queries"

	"github.com/jackc/pgx/v5/pgxpool"
)

func FindSession(
	ctx context.Context,
	pool *pgxpool.Pool,
	tokenHash [32]byte,
) (model.Session, error) {
	session := model.Session{
		TokenHash: tokenHash,
	}

	err := pool.QueryRow(
		ctx, queries.FindSessionUser, session.TokenHash[:],
	).Scan(
		&session.User.ID,
		&session.User.Username,
		&session.User.Email,
		&session.User.ProfileImage,
		&session.User.BodyWeight,
		&session.User.Height,
		&session.User.Gender,
		&session.User.DateOfBirth,
		&session.User.CreatedAt,
		&session.User.UpdatedAt,
	)

	return session, err
}

func DeleteSession(
	ctx context.Context,
	pool *pgxpool.Pool,
	tokenHash [32]byte,
) error {
	_, err := pool.Exec(ctx, queries.DeleteSession, tokenHash[:])

	return err
}

func CreateSession(
	ctx context.Context,
	pool *pgxpool.Pool,
	tokenHash [32]byte,
	userID int,
	expiresAt time.Time,
) error {
	_, err := pool.Exec(
		ctx, queries.CreateSession, tokenHash[:], userID, expiresAt,
	)

	return err
}
