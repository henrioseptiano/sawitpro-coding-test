// This file contains types that are used in the repository layer.
package repository

import "time"

type User struct {
	ID                       int        `json:"id" gorm:"id,primaryKey,autoIncrement"`
	UserID                   string     `json:"user_id" gorm:"user_id,unique"`
	FullName                 string     `json:"full_name" gorm:"full_name,not null"`
	PhoneNumber              string     `json:"phone_number" gorm:"phone_number,not null"`
	Password                 string     `json:"password" gorm:"password,not null"`
	SuccessfullLoginAttempts int64      `json:"sucessfull_login_attempts" gorm:"successfull_login_attempts, not null"`
	LastLogin                *time.Time `json:"last_login" gorm:"last_login"`
	CreatedAt                time.Time  `json:"created_at" gorm:"created_at,not null"`
	UpdatedAt                time.Time  `json:"updated_at" gorm:"updated_at,not null"`
}

type GetTestByIdInput struct {
	Id string
}

type GetTestByIdOutput struct {
	Name string
}
