package repository

import (
	"context"
	"encoding/json"
	"errors"

	"kheun-Lift-backend-go/internal/leaderboards/model"
	"kheun-Lift-backend-go/internal/queries"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrExerciseNotFound    = errors.New("Exercise not found")
	ErrWeightClassNotFound = errors.New("Weight Class not found")
	ErrWeightClassGender   = errors.New("Weight class does not match gender")
)

func ValidateFilters(ctx context.Context, pool *pgxpool.Pool, exerciseID int32, weightClassID *int32, gender *string) error {
	var exists bool
	if err := pool.QueryRow(ctx, queries.LiftExerciseExists, exerciseID).Scan(&exists); err != nil {
		return err
	}
	if !exists {
		return ErrExerciseNotFound
	}
	if weightClassID == nil {
		return nil
	}

	var classGender string
	err := pool.QueryRow(ctx, queries.LeaderboardWeightClassGender, *weightClassID).Scan(&classGender)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrWeightClassNotFound
	}
	if err != nil {
		return err
	}
	if gender != nil && *gender != classGender {
		return ErrWeightClassGender
	}
	return nil
}

func List(ctx context.Context, pool *pgxpool.Pool, exerciseID int32, weightClassID *int32, gender *string) ([]json.RawMessage, error) {
	rows, err := pool.Query(ctx, queries.ListLeaderboard, exerciseID, weightClassID, gender)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowTo[json.RawMessage])
}

func Random(ctx context.Context, pool *pgxpool.Pool) (*model.RandomLeaderboard, error) {
	var result model.RandomLeaderboard
	err := pool.QueryRow(ctx, queries.RandomLeaderboardPair).Scan(&result.ExerciseID, &result.WeightClassID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	result.Leaderboards, err = List(ctx, pool, result.ExerciseID, &result.WeightClassID, nil)
	if err != nil {
		return nil, err
	}
	if len(result.Leaderboards) == 0 {
		return nil, nil
	}
	return &result, nil
}
