package model

import "encoding/json"

type RandomLeaderboard struct {
	ExerciseID    int32             `json:"exerciseId"`
	WeightClassID int32             `json:"weightClassId"`
	Leaderboards  []json.RawMessage `json:"leaderboards"`
}
