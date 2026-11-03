package lifts

import (
	"context"
	"errors"

	"kheun-Lift-backend-go/internal/queries"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrDuplicate   = errors.New("This exercise, weight and reps were already posted.")
	ErrProfile     = errors.New("Body weight and gender are required.")
	ErrExercise    = errors.New("Exercise not found.")
	ErrWeightClass = errors.New("Weight class not found.")
)

func scanLift(row pgx.Row) (LiftRecord, error) {
	var record LiftRecord
	err := row.Scan(
		&record.ID, &record.UserID, &record.ExerciseID, &record.WeightClassID,
		&record.Weight, &record.Reps, &record.Caption, &record.VideoURL,
		&record.Status, &record.CreatedAt, &record.UpdatedAt,
		&record.User.Username, &record.Exercise.Name,
		&record.WeightClass.Name, &record.WeightClass.Gender,
	)
	record.Exercise.ID = record.ExerciseID
	record.WeightClass.ID = record.WeightClassID
	return record, err
}

func listByUser(ctx context.Context, pool *pgxpool.Pool, userID int) ([]LiftRecord, error) {
	rows, err := pool.Query(ctx, queries.ListLiftsByUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := make([]LiftRecord, 0)
	for rows.Next() {
		record, err := scanLift(rows)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, rows.Err()
}

func getByID(ctx context.Context, pool *pgxpool.Pool, id int32) (LiftRecord, error) {
	return scanLift(pool.QueryRow(ctx, queries.GetLiftByID, id))
}

func create(ctx context.Context, pool *pgxpool.Pool, userID int, input CreateInput) (LiftRecord, error) {
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return LiftRecord{}, err
	}
	defer tx.Rollback(ctx)

	var bodyWeight *float64
	var gender *string
	// shortcut: serialize Go creates per user; add a unique constraint before allowing other writers.
	if err := tx.QueryRow(ctx, queries.LockLiftUser, userID).Scan(&bodyWeight, &gender); err != nil {
		return LiftRecord{}, err
	}
	if bodyWeight == nil || gender == nil {
		return LiftRecord{}, ErrProfile
	}

	var exists bool
	if err := tx.QueryRow(ctx, queries.LiftExerciseExists, input.ExerciseID).Scan(&exists); err != nil {
		return LiftRecord{}, err
	}
	if !exists {
		return LiftRecord{}, ErrExercise
	}
	if err := tx.QueryRow(ctx, queries.LiftExists, userID, input.ExerciseID, input.Weight, input.Reps).Scan(&exists); err != nil {
		return LiftRecord{}, err
	}
	if exists {
		return LiftRecord{}, ErrDuplicate
	}

	var weightClassID int
	err = tx.QueryRow(ctx, queries.FindLiftWeightClass, *gender, *bodyWeight).Scan(&weightClassID)
	if errors.Is(err, pgx.ErrNoRows) {
		return LiftRecord{}, ErrWeightClass
	}
	if err != nil {
		return LiftRecord{}, err
	}

	record, err := scanLift(tx.QueryRow(ctx, queries.InsertLift,
		userID, input.ExerciseID, weightClassID, input.Weight, input.Reps, input.Caption, input.VideoURL,
	))
	if err != nil {
		return LiftRecord{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return LiftRecord{}, err
	}
	return record, nil
}

func updateCaption(ctx context.Context, pool *pgxpool.Pool, id int32, userID int, caption string) (LiftRecord, error) {
	return scanLift(pool.QueryRow(ctx, queries.UpdateLiftCaption, id, userID, caption))
}

func deleteByID(ctx context.Context, pool *pgxpool.Pool, id int32, userID int) error {
	result, err := pool.Exec(ctx, queries.DeleteLift, id, userID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
