package model

import "time"

type RegisterInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID           int        `json:"id"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	ProfileImage *string    `json:"profileImage"`
	BodyWeight   *float64   `json:"bodyWeight"`
	Height       *float64   `json:"height"`
	Gender       *string    `json:"gender"`
	DateOfBirth  *time.Time `json:"dateOfBirth"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}

type LoginUser struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type LoginResult struct {
	User      LoginUser
	Token     string `json:"-"`
	ExpiresAt time.Time
}

type Session struct {
	User      User
	TokenHash [32]byte `json:"-"`
}

type Exercise struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}
