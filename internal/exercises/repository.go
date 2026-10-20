package exercises

import (
	"context"

	"kheun-Lift-backend-go/internal/queries"

	"github.com/jackc/pgx/v5/pgxpool"
)

func list(ctx context.Context, pool *pgxpool.Pool) ([]Exercise, error) {
	rows, err := pool.Query(ctx, queries.ListExercises)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	exercises := make([]Exercise, 0)

	for rows.Next() {
		var exercise Exercise

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
