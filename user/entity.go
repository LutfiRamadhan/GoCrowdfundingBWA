package user

import "time"

type User struct {
	ID         string `json:"id"`
	Name       string
	Occupation string
	Email      string
	Password   string
	Profilepic string
	Role       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
