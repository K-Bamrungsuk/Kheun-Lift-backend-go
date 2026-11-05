package model

import "time"

type CreateInput struct {
	ExerciseID int32   `json:"exerciseId"`
	Weight     float64 `json:"weight"`
	Reps       int32   `json:"reps"`
	Caption    *string `json:"caption"`
	VideoURL   *string `json:"videoUrl"`
}

type LiftRecord struct {
	ID            int       `json:"id"`
	UserID        int       `json:"userId"`
	ExerciseID    int       `json:"exerciseId"`
	WeightClassID int       `json:"weightClassId"`
	Weight        float64   `json:"weight"`
	Reps          int       `json:"reps"`
	Caption       *string   `json:"caption"`
	VideoURL      *string   `json:"videoUrl"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	User          struct {
		Username string `json:"username"`
	} `json:"user"`
	Exercise struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"exercise"`
	WeightClass struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Gender string `json:"gender"`
	} `json:"weightClass"`
}

type UpdateInput struct {
	Caption *string `json:"caption"`
}
