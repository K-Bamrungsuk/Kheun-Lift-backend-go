package repository

import (
	"context"

	"kheun-Lift-backend-go/internal/exercises/model"
	"kheun-Lift-backend-go/internal/queries"

	"github.com/jackc/pgx/v5/pgxpool"
)

func List(ctx context.Context, pool *pgxpool.Pool) ([]model.Exercise, error) {
	rows, err := pool.Query(ctx, queries.ListExercises)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	exercises := make([]model.Exercise, 0)

	for rows.Next() {
		var exercise model.Exercise

		if err := rows.Scan(
			&exercise.ID,
			&exercise.Name,
			&exercise.Description,
		); err != nil {
			return nil, err
		}

		exercises = append(exercises, exercise)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return exercises, nil
}

func GetByID(ctx context.Context, pool *pgxpool.Pool, id int32) (model.Exercise, error) {
	var exercise model.Exercise

	err := pool.QueryRow(ctx, queries.GetExerciseByID, id).Scan(
		&exercise.ID,
		&exercise.Name,
		&exercise.Description,
	)

	return exercise, err
}
