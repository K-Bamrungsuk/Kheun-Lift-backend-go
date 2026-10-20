package model

import "time"

type UpdateUserInput struct {
	Username    *string
	BodyWeight  *float64
	Height      *float64
	Gender      *string
	DateOfBirth *time.Time
}
