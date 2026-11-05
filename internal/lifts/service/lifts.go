package service

import (
	"context"
	"errors"

	"kheun-Lift-backend-go/internal/lifts/model"
	"kheun-Lift-backend-go/internal/lifts/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrDuplicate   = errors.New("This exercise, weight and reps were already posted.")
	ErrProfile     = errors.New("Body weight and gender are required.")
	ErrExercise    = errors.New("Exercise not found.")
	ErrWeightClass = errors.New("Weight class not found.")
)

func Create(ctx context.Context, pool *pgxpool.Pool, userID int, input model.CreateInput) (model.LiftRecord, error) {
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return model.LiftRecord{}, err
	}
	defer tx.Rollback(ctx)

	bodyWeight, gender, err := repository.LockUser(ctx, tx, userID)
	if err != nil {
		return model.LiftRecord{}, err
	}
	if bodyWeight == nil || gender == nil {
		return model.LiftRecord{}, ErrProfile
	}

	exists, err := repository.ExerciseExists(ctx, tx, input.ExerciseID)
	if err != nil {
		return model.LiftRecord{}, err
	}
	if !exists {
		return model.LiftRecord{}, ErrExercise
	}
	exists, err = repository.Exists(ctx, tx, userID, input)
	if err != nil {
		return model.LiftRecord{}, err
	}
	if exists {
		return model.LiftRecord{}, ErrDuplicate
	}

	weightClassID, err := repository.FindWeightClassID(ctx, tx, *gender, *bodyWeight)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.LiftRecord{}, ErrWeightClass
	}
	if err != nil {
		return model.LiftRecord{}, err
	}

	record, err := repository.Insert(ctx, tx, userID, weightClassID, input)
	if err != nil {
		return model.LiftRecord{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return model.LiftRecord{}, err
	}
	return record, nil
}
