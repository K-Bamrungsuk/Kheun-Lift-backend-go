package model

import "time"

type UpdateUserInput struct {
	Username     *string
	ProfileImage *string
	BodyWeight   *float64
	Height       *float64
	Gender       *string
	DateOfBirth  *time.Time
}
