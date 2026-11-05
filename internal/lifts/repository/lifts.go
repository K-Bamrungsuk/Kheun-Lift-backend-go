package repository

import (
	"context"

	"kheun-Lift-backend-go/internal/lifts/model"
	"kheun-Lift-backend-go/internal/queries"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func scanLift(row pgx.Row) (model.LiftRecord, error) {
	var record model.LiftRecord
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

func ListByUser(ctx context.Context, pool *pgxpool.Pool, userID int) ([]model.LiftRecord, error) {
	rows, err := pool.Query(ctx, queries.ListLiftsByUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := make([]model.LiftRecord, 0)
	for rows.Next() {
		record, err := scanLift(rows)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, rows.Err()
}

func GetByID(ctx context.Context, pool *pgxpool.Pool, id int32) (model.LiftRecord, error) {
	return scanLift(pool.QueryRow(ctx, queries.GetLiftByID, id))
}

func UpdateCaption(ctx context.Context, pool *pgxpool.Pool, id int32, userID int, caption string) (model.LiftRecord, error) {
	return scanLift(pool.QueryRow(ctx, queries.UpdateLiftCaption, id, userID, caption))
}

func DeleteByID(ctx context.Context, pool *pgxpool.Pool, id int32, userID int) error {
	result, err := pool.Exec(ctx, queries.DeleteLift, id, userID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func LockUser(ctx context.Context, tx pgx.Tx, userID int) (*float64, *string, error) {
	var bodyWeight *float64
	var gender *string
	// shortcut: serialize Go creates per user; add a unique constraint before allowing other writers.
	err := tx.QueryRow(ctx, queries.LockLiftUser, userID).Scan(&bodyWeight, &gender)
	return bodyWeight, gender, err
}

func ExerciseExists(ctx context.Context, tx pgx.Tx, exerciseID int32) (bool, error) {
	var exists bool
	err := tx.QueryRow(ctx, queries.LiftExerciseExists, exerciseID).Scan(&exists)
	return exists, err
}

func Exists(ctx context.Context, tx pgx.Tx, userID int, input model.CreateInput) (bool, error) {
	var exists bool
	err := tx.QueryRow(ctx, queries.LiftExists, userID, input.ExerciseID, input.Weight, input.Reps).Scan(&exists)
	return exists, err
}

func FindWeightClassID(ctx context.Context, tx pgx.Tx, gender string, bodyWeight float64) (int, error) {
	var id int
	err := tx.QueryRow(ctx, queries.FindLiftWeightClass, gender, bodyWeight).Scan(&id)
	return id, err
}

func Insert(ctx context.Context, tx pgx.Tx, userID, weightClassID int, input model.CreateInput) (model.LiftRecord, error) {
	return scanLift(tx.QueryRow(ctx, queries.InsertLift,
		userID, input.ExerciseID, weightClassID, input.Weight, input.Reps, input.Caption, input.VideoURL,
	))
}
