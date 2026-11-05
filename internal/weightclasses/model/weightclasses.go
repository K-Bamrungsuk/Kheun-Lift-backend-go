package model

import "time"

type WeightClass struct {
	ID        int       `json:"id"`
	Gender    string    `json:"gender"`
	Name      string    `json:"name"`
	MinWeight float64   `json:"minWeight"`
	MaxWeight *float64  `json:"maxWeight"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
