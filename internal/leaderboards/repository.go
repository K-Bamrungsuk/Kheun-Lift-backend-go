package leaderboards

import (
	"context"
	"encoding/json"
	"errors"

	"kheun-Lift-backend-go/internal/queries"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	errExerciseNotFound    = errors.New("Exercise not found")
	errWeightClassNotFound = errors.New("Weight Class not found")
	errWeightClassGender   = errors.New("Weight class does not match gender")
)

type randomLeaderboard struct {
	ExerciseID    int32             `json:"exerciseId"`
	WeightClassID int32             `json:"weightClassId"`
	Leaderboards  []json.RawMessage `json:"leaderboards"`
}

func validateFilters(ctx context.Context, pool *pgxpool.Pool, exerciseID int32, weightClassID *int32, gender *string) error {
	var exists bool
	if err := pool.QueryRow(ctx, queries.LiftExerciseExists, exerciseID).Scan(&exists); err != nil {
		return err
	}
	if !exists {
		return errExerciseNotFound
	}
	if weightClassID == nil {
		return nil
	}

	var classGender string
	err := pool.QueryRow(ctx, queries.LeaderboardWeightClassGender, *weightClassID).Scan(&classGender)
	if errors.Is(err, pgx.ErrNoRows) {
		return errWeightClassNotFound
	}
	if err != nil {
		return err
	}
	if gender != nil && *gender != classGender {
		return errWeightClassGender
	}
	return nil
}

func list(ctx context.Context, pool *pgxpool.Pool, exerciseID int32, weightClassID *int32, gender *string) ([]json.RawMessage, error) {
	rows, err := pool.Query(ctx, queries.ListLeaderboard, exerciseID, weightClassID, gender)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowTo[json.RawMessage])
}

func random(ctx context.Context, pool *pgxpool.Pool) (*randomLeaderboard, error) {
	var result randomLeaderboard
	err := pool.QueryRow(ctx, queries.RandomLeaderboardPair).Scan(&result.ExerciseID, &result.WeightClassID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	result.Leaderboards, err = list(ctx, pool, result.ExerciseID, &result.WeightClassID, nil)
	if err != nil {
		return nil, err
	}
	if len(result.Leaderboards) == 0 {
		return nil, nil
	}
	return &result, nil
}
