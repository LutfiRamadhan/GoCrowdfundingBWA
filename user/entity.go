package user

import "time"

type User struct {
	ID         int    `json:"id" gorm:"id,primaryKey,autoIncrement"`
	Name       string `json:"name" gorm:"name,index"`
	Occupation string `json:"occupation" gorm:"occupation,index"`
	Email      string `json:"email" gorm:"email,uniqueIndex"`
	Password   string `json:"password" gorm:"password"`
	ProfilePic string `json:"profile_pic" gorm:"profile_pic"`
	Role       string `json:"role" gorm:"role"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
