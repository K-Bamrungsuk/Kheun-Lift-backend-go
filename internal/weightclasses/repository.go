package weightclasses

import (
	"context"
	"errors"

	"kheun-Lift-backend-go/internal/queries"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func FindForUser(ctx context.Context, pool *pgxpool.Pool, gender *string, bodyWeight *float64) (*WeightClass, error) {
	if gender == nil || bodyWeight == nil {
		return nil, nil
	}
	var weightClass WeightClass
	err := pool.QueryRow(ctx, queries.FindUserWeightClass, *gender, *bodyWeight).Scan(
		&weightClass.ID, &weightClass.Gender, &weightClass.Name,
		&weightClass.MinWeight, &weightClass.MaxWeight,
		&weightClass.CreatedAt, &weightClass.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &weightClass, nil
}

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
