package weightclasses

import (
	"context"

	"kheun-Lift-backend-go/internal/queries"

	"github.com/jackc/pgx/v5/pgxpool"
)

func list(ctx context.Context, pool *pgxpool.Pool, gender string) ([]WeightClass, error) {
	rows, err := pool.Query(ctx, queries.ListWeightClasses, gender)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	weightClasses := make([]WeightClass, 0)

	for rows.Next() {
		var weightClass WeightClass

		if err := rows.Scan(
			&weightClass.ID,
			&weightClass.Gender,
			&weightClass.Name,
			&weightClass.MinWeight,
			&weightClass.MaxWeight,
			&weightClass.CreatedAt,
			&weightClass.UpdatedAt,
		); err != nil {
			return nil, err
		}

		weightClasses = append(weightClasses, weightClass)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return weightClasses, nil
}
