package model

import (
	"time"

	weightclassmodel "kheun-Lift-backend-go/internal/weightclasses/model"
)

type PublicUser struct {
	ID           int                           `json:"id"`
	Username     string                        `json:"username"`
	ProfileImage *string                       `json:"profileImage"`
	Gender       *string                       `json:"gender"`
	BodyWeight   *float64                      `json:"bodyWeight"`
	Height       *float64                      `json:"height"`
	WeightClass  *weightclassmodel.WeightClass `json:"weightClass"`
}

type UpdateUserInput struct {
	Username     *string
	ProfileImage *string
	BodyWeight   *float64
	Height       *float64
	Gender       *string
	DateOfBirth  *time.Time
}
